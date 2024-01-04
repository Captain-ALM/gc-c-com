package transport

import (
	"bufio"
	"errors"
	"golang.local/gc-c-com/packet"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Handler struct {
	ID            string
	active        bool
	closeEvent    func(t Transport, e error)
	sendMutex     *sync.Mutex
	sendBuffer    [][]byte
	recvMutex     *sync.Mutex
	recvNotif     chan [][]byte
	recvBuffer    [][]byte
	connNotif     chan bool
	timeout       time.Duration
	closeMutex    *sync.Mutex
	closedChannel chan bool
}

func (h *Handler) Activate() {
	if h == nil || h.IsActive() {
		return
	}
	h.closedChannel = make(chan bool)
	h.sendMutex = &sync.Mutex{}
	h.recvMutex = &sync.Mutex{}
	h.closeMutex = &sync.Mutex{}
	h.sendBuffer = nil
	h.recvBuffer = nil
	h.recvNotif = make(chan [][]byte)
	h.connNotif = make(chan bool)
	h.active = true
	go func() {
		ctOut := h.timeout
		tOut := time.NewTimer(ctOut)
		defer func() {
			if !tOut.Stop() {
				select {
				case <-tOut.C:
				default:
				}
			}
		}()
		for h.active {
			select {
			case <-h.closedChannel:
				return
			case valid := <-h.connNotif:
				ctOut = h.timeout
				if valid {
					if !tOut.Stop() {
						select {
						case <-tOut.C:
						default:
						}
					}
					tOut.Reset(ctOut)
				} else {
					valid = true
					for valid {
						select {
						case <-h.closedChannel:
							return
						case <-h.connNotif:
						case <-tOut.C:
							tOut.Reset(ctOut)
							valid = false
						}
					}
				}
			case <-tOut.C:
				h.closeMutex.Lock()
				defer h.closeMutex.Unlock()
				if h.active {
					h.active = false
					_ = h.close(errors.New("handler timeout"))
				}
				return
			}
		}
	}()
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !h.IsActive() {
		return
	}
	select {
	case <-h.closedChannel:
		return
	case h.connNotif <- true:
	default:
	}
	defer func() {
		select {
		case <-h.closedChannel:
			return
		case h.connNotif <- true:
		default:
		}
	}()
	hasPing := false
	if request.Method != http.MethodOptions {
		hasPing = h.receiveRequest(request)
	}
	if request.Method == http.MethodGet || request.Method == http.MethodPost {
		h.sendResponse(writer, hasPing)
	} else if request.Method == http.MethodDelete {
		h.closeMutex.Lock()
		defer h.closeMutex.Unlock()
		if h.active {
			h.active = false
			_ = h.close(errors.New("handler closed remotely"))
		}
		writer.WriteHeader(http.StatusAccepted)
	} else if request.Method == http.MethodOptions {
		writer.Header().Set("Allow", http.MethodOptions+", "+http.MethodGet+", "+http.MethodPost+", "+http.MethodDelete)
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) receiveRequest(request *http.Request) bool {
	if request.Body == nil {
		return false
	}
	hasPing := false
	bScan := bufio.NewScanner(request.Body)
	var rIn [][]byte
	for bScan.Scan() {
		cBts := bScan.Bytes()
		cR := make([]byte, len(cBts))
		copy(cR, cBts)
		switch packet.GetCommandIgnoreError(cR) {
		case packet.Ping:
			hasPing = true
		case packet.Pong:
		default:
			rIn = append(rIn, cR)
		}
	}
	select {
	case <-h.closedChannel:
		return hasPing
	case h.recvNotif <- rIn:
	}
	return hasPing
}

func (h *Handler) sendResponse(response http.ResponseWriter, needPong bool) {
	if response == nil {
		return
	}
	h.sendMutex.Lock()
	defer h.sendMutex.Unlock()
	sz := h.getSendSize()
	var thePong []byte
	if needPong {
		thePong = packet.NewPong().ToBytesIgnoreError()
		sz += len(thePong) + 2
	}
	if sz < 1 {
		response.Header().Set("Content-Length", "0")
		response.WriteHeader(http.StatusAccepted)
	} else {
		response.Header().Set("Content-Length", strconv.Itoa(sz))
		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		defer func() { h.sendBuffer = nil }()
		if len(thePong) > 0 {
			h.sendBuffer = append(h.sendBuffer, thePong)
		}
		for _, bytes := range h.sendBuffer {
			_, err := response.Write(bytes)
			if err != nil {
				h.closeMutex.Lock()
				defer h.closeMutex.Unlock()
				if h.active {
					h.active = false
					_ = h.close(err)
				}
				return
			}
			_, err = response.Write([]byte("\r\n"))
			if err != nil {
				h.closeMutex.Lock()
				defer h.closeMutex.Unlock()
				if h.active {
					h.active = false
					_ = h.close(err)
				}
				return
			}
		}
	}
}

func (h *Handler) getSendSize() int {
	sz := 0
	for _, bytes := range h.sendBuffer {
		sz += len(bytes) + 2
	}
	return sz
}

func (h *Handler) GetID() string {
	if h == nil {
		return ""
	}
	return h.ID
}

func (h *Handler) IsActive() bool {
	if h == nil {
		return false
	}
	return h.active
}

func (h *Handler) Send(p *packet.Packet) error {
	if !h.IsActive() {
		return errors.New("handler not active")
	}
	bts, err := p.ToBytes()
	if err != nil {
		return err
	}
	h.sendMutex.Lock()
	defer h.sendMutex.Unlock()
	h.sendBuffer = append(h.sendBuffer, bts)
	return nil
}

func (h *Handler) Receive() (*packet.Packet, error) {
	if !h.IsActive() {
		return nil, errors.New("handler not active")
	}
	h.recvMutex.Lock()
	defer h.recvMutex.Unlock()
	if len(h.recvBuffer) < 1 {
		select {
		case <-h.closedChannel:
			return nil, errors.New("handler closed")
		case pl := <-h.recvNotif:
			h.recvBuffer = append(h.recvBuffer, pl...)
		}
	}
	if len(h.recvBuffer) > 0 {
		var tr []byte
		tr, h.recvBuffer = h.recvBuffer[0], h.recvBuffer[1:]
		return packet.From(tr)
	}
	return nil, errors.New("handler receive empty")
}

func (h *Handler) close(err error) error {
	close(h.closedChannel)
	//close(h.recvNotif)
	//close(h.connNotif)
	if h.closeEvent != nil {
		h.closeEvent(h, err)
	}
	h.recvBuffer = nil
	h.sendBuffer = nil
	return err
}

func (h *Handler) Close() error {
	if h == nil {
		return nil
	}
	h.closeMutex.Lock()
	defer h.closeMutex.Unlock()
	if h.active {
		h.active = false
		return h.close(nil)
	}
	return nil
}

func (h *Handler) SetOnClose(callback func(t Transport, e error)) {
	if h == nil || callback == nil {
		return
	}
	h.closeEvent = callback
}

func (h *Handler) SetTimeout(to time.Duration) {
	if h == nil {
		return
	}
	h.timeout = to
	if h.active {
		select {
		case <-h.closedChannel:
			return
		case h.connNotif <- false:
		}
	}
}

func (h *Handler) GetTimeout() time.Duration {
	if h == nil {
		return 0
	}
	return h.timeout
}
