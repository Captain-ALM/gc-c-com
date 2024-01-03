package transport

import (
	"errors"
	"github.com/gorilla/websocket"
	"golang.local/gc-c-com/packet"
	"sync"
	"time"
)

type Websocket struct {
	ID            string
	active        bool
	closeEvent    func(t Transport, e error)
	conn          *websocket.Conn
	recvNotif     chan []byte
	sendNotif     chan []byte
	timeout       time.Duration
	closeMutex    *sync.Mutex
	closedChannel chan bool
}

func (w *Websocket) Activate(conn *websocket.Conn) {
	if w == nil || conn == nil || w.IsActive() {
		return
	}
	w.conn = conn
	w.closedChannel = make(chan bool)
	w.closeMutex = &sync.Mutex{}
	w.recvNotif = make(chan []byte)
	w.sendNotif = make(chan []byte)
	w.active = true
	go func() {
		var recvBuff []byte
		defer func() { _ = w.conn.Close() }()
		for w.active {
			err := w.conn.SetReadDeadline(time.Now().Add(w.timeout))
			if err != nil {
				w.closeMutex.Lock()
				defer w.closeMutex.Unlock()
				if w.active {
					w.active = false
					_ = w.close(err)
				}
				return
			}
			_, msg, err := w.conn.ReadMessage()
			if err != nil {
				w.closeMutex.Lock()
				defer w.closeMutex.Unlock()
				if w.active {
					w.active = false
					_ = w.close(err)
				}
				return
			}
			if len(recvBuff) < 1 && msg[len(msg)-1] == '\n' {
				switch packet.GetCommandIgnoreError(msg) {
				case packet.Ping:
					select {
					case <-w.closedChannel:
					case w.sendNotif <- packet.NewPong().ToBytesIgnoreError():
					}
				case packet.Pong:
				default:
					select {
					case <-w.closedChannel:
					case w.recvNotif <- msg:
					}
				}
			} else if len(recvBuff) > 0 && msg[len(msg)-1] == '\n' {
				recvBuff = append(recvBuff, msg...)
				switch packet.GetCommandIgnoreError(recvBuff) {
				case packet.Ping:
					select {
					case <-w.closedChannel:
					case w.sendNotif <- packet.NewPong().ToBytesIgnoreError():
					}
				case packet.Pong:
				default:
					select {
					case <-w.closedChannel:
					case w.recvNotif <- recvBuff:
					}
				}
				recvBuff = nil
			} else {
				recvBuff = append(recvBuff, msg...)
			}
		}
	}()
	go func() {
		defer func() { _ = w.conn.Close() }()
		for w.active {
			err := w.conn.SetWriteDeadline(time.Now().Add(w.timeout))
			if err != nil {
				w.closeMutex.Lock()
				defer w.closeMutex.Unlock()
				if w.active {
					w.active = false
					_ = w.close(err)
				}
				return
			}
			select {
			case <-w.closedChannel:
				_ = w.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			case msg := <-w.sendNotif:
				err = w.conn.WriteMessage(websocket.TextMessage, append(msg, []byte("\r\n")...))
				if err != nil {
					w.closeMutex.Lock()
					defer w.closeMutex.Unlock()
					if w.active {
						w.active = false
						_ = w.close(err)
					}
					return
				}
			}
		}
	}()
}

func (w *Websocket) GetID() string {
	if w == nil {
		return ""
	}
	return w.ID
}

func (w *Websocket) IsActive() bool {
	if w == nil {
		return false
	}
	return w.active
}

func (w *Websocket) Send(p *packet.Packet) error {
	if !w.IsActive() {
		return errors.New("websocket not active")
	}
	bts, err := p.ToBytes()
	if err != nil {
		return err
	}
	select {
	case <-w.closedChannel:
	case w.sendNotif <- bts:
	}
	return nil
}

func (w *Websocket) Receive() (*packet.Packet, error) {
	if !w.IsActive() {
		return nil, errors.New("websocket not active")
	}
	select {
	case <-w.closedChannel:
		return nil, errors.New("websocket not active")
	case msg := <-w.recvNotif:
		return packet.From(msg)
	}
}

func (w *Websocket) close(err error) error {
	_ = w.conn.Close()
	close(w.closedChannel)
	//close(w.recvNotif)
	//close(w.sendNotif)
	if w.closeEvent != nil {
		w.closeEvent(w, err)
	}
	return err
}

func (w *Websocket) Close() error {
	if w == nil {
		return nil
	}
	w.closeMutex.Lock()
	defer w.closeMutex.Unlock()
	if w.active {
		w.active = false
		_ = w.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(w.timeout))
		return w.close(nil)
	}
	return nil
}

func (w *Websocket) SetOnClose(callback func(t Transport, e error)) {
	if w == nil || callback == nil {
		return
	}
	w.closeEvent = callback
}

func (w *Websocket) SetTimeout(to time.Duration) {
	if w == nil {
		return
	}
	w.timeout = to
}

func (w *Websocket) GetTimeout() time.Duration {
	if w == nil {
		return 0
	}
	return w.timeout
}
