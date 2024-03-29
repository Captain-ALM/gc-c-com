package transport

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/Captain-ALM/gc-c-com/packet"
	"github.com/gorilla/websocket"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	restTargetURL string
	restClient    *http.Client
	closeNotif    chan bool
	recvNotif     chan [][]byte
	sendNotif     chan []byte
	closeMutex    *sync.Mutex
	sendMutex     *sync.Mutex
	sendBuffer    [][]byte
	recvMutex     *sync.Mutex
	recvBuffer    [][]byte
	conn          *websocket.Conn
	wsDialer      *websocket.Dialer
	keepAlive     time.Duration
	kaMutex       *sync.Mutex
	kaNotif       chan bool
	timeout       time.Duration
	active        bool
	closeEvent    func(t Transport, e error)
}

func (c *Client) Activate(wsURL string, restURL string) {
	if c == nil || c.IsActive() {
		return
	}
	c.closeNotif = make(chan bool)
	c.recvNotif = make(chan [][]byte)
	c.sendNotif = make(chan []byte)
	c.closeMutex = &sync.Mutex{}
	c.sendMutex = &sync.Mutex{}
	c.recvMutex = &sync.Mutex{}
	c.kaMutex = &sync.Mutex{}
	c.kaNotif = make(chan bool)
	c.active = true
	if wsURL != "" {
		c.wsStart(wsURL)
		if c.conn != nil {
			return
		}
	}
	c.active = c.restStart(restURL)
}

func (c *Client) restStart(restURL string) bool {
	c.restClient = &http.Client{Timeout: c.timeout}
	resp, err := c.restClient.Get(restURL)
	if err != nil {
		debugErrIsNil(err)
		return false
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusOK && resp.ContentLength > 0 && strings.HasPrefix(strings.ToLower(resp.Header.Get("Content-Type")), "text/plain") {
		buff := make([]byte, resp.ContentLength)
		ln, err := resp.Body.Read(buff)
		if err != nil {
			debugErrIsNil(err)
			return false
		}
		c.restTargetURL = restURL + "?s=" + url.QueryEscape(strings.Trim(string(buff[:ln]), "\r\n"))
		if c.keepAlive < 1 {
			c.keepAlive = time.Second
		}
		go func() {
			for c.active {
				select {
				case <-c.closeNotif:
					return
				case bts := <-c.sendNotif:
					tLn := c.appendToSendBuffer(bts)
					debugPrintln("Client.restStart PUMP - tLn: " + strconv.Itoa(tLn))
				}
			}
		}()
		go func() {
			kAlive := time.NewTicker(c.keepAlive)
			defer kAlive.Stop()
			for c.active {
				fail, rPong := c.handlerProcessor()
				if fail {
					c.closeMutex.Lock()
					defer c.closeMutex.Unlock()
					if c.active {
						c.active = false
						_ = c.close(nil)
						return
					}
				}
				if rPong {
					select {
					case <-c.closeNotif:
					case c.sendNotif <- packet.NewPong().ToBytesIgnoreError():
					}
				}
				select {
				case <-kAlive.C:
					kAlive.Reset(c.keepAlive)
				case <-c.closeNotif:
					c.closeMutex.Lock()
					defer c.closeMutex.Unlock()
					if c.active {
						c.active = false
						_ = c.close(nil)
						return
					}
				}
			}
		}()
		return true
	}
	return false
}

func (c *Client) appendToSendBuffer(bts []byte) int {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()
	c.sendBuffer = append(c.sendBuffer, bts)
	return len(c.sendBuffer)
}

func (c *Client) handlerProcessor() (failed bool, hasPing bool) {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()
	var resp *http.Response
	var err error
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		c.sendBuffer = nil
	}()
	if len(c.sendBuffer) < 1 {
		resp, err = c.restClient.Get(c.restTargetURL)
	} else {
		buff := bytes.NewBuffer(make([]byte, 0, c.getSendSize()))
		for _, bts := range c.sendBuffer {
			_, err = buff.Write(bts)
			if err != nil {
				debugErrIsNil(err)
				return true, false
			}
			_, err = buff.Write([]byte("\r\n"))
			if err != nil {
				debugErrIsNil(err)
				return true, false
			}
		}
		debugPrintln("Client.handlerProcessor - buff: " + hex.EncodeToString(buff.Bytes()))
		resp, err = c.restClient.Post(c.restTargetURL, "text/plain; charset=utf-8", buff)
	}
	if err != nil {
		debugErrIsNil(err)
		return true, false
	}
	debugPrintln("Client.handlerProcessor - cl: " + strconv.Itoa(int(resp.ContentLength)))
	if resp.StatusCode == http.StatusOK && resp.ContentLength > 0 && strings.HasPrefix(strings.ToLower(resp.Header.Get("Content-Type")), "text/plain") {
		hasPing = false
		bScan := bufio.NewScanner(resp.Body)
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
		case <-c.closeNotif:
			return true, hasPing
		case c.recvNotif <- rIn:
			debugPrintln("Client.handlerProcessor - rl: " + strconv.Itoa(len(rIn)))
		}
		return false, hasPing
	} else if resp.StatusCode == http.StatusAccepted {
		return false, false
	}
	return true, false
}

func (c *Client) getSendSize() int {
	sz := 0
	for _, bts := range c.sendBuffer {
		sz += len(bts) + 2
	}
	return sz
}

func (c *Client) wsStart(wsURL string) {
	c.conn = nil
	c.wsDialer = &websocket.Dialer{HandshakeTimeout: c.timeout}
	conn, _, err := c.wsDialer.Dial(wsURL, nil)
	if err != nil {
		debugErrIsNil(err)
		return
	}
	c.conn = conn
	go func() {
		var recvBuff []byte
		defer func() { _ = c.conn.Close() }()
		for c.active {
			err := c.conn.SetReadDeadline(time.Now().Add(c.timeout))
			if err != nil {
				c.closeMutex.Lock()
				defer c.closeMutex.Unlock()
				if c.active {
					c.active = false
					_ = c.close(err)
				}
				return
			}
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				c.closeMutex.Lock()
				defer c.closeMutex.Unlock()
				if c.active {
					c.active = false
					_ = c.close(err)
				}
				return
			}
			if len(recvBuff) < 1 && msg[len(msg)-1] == '\n' {
				switch packet.GetCommandIgnoreError(msg) {
				case packet.Ping:
					select {
					case <-c.closeNotif:
					case c.sendNotif <- packet.NewPong().ToBytesIgnoreError():
					}
				case packet.Pong:
				default:
					select {
					case <-c.closeNotif:
					case c.recvNotif <- [][]byte{msg}:
					}
				}
			} else if len(recvBuff) > 0 && msg[len(msg)-1] == '\n' {
				recvBuff = append(recvBuff, msg...)
				switch packet.GetCommandIgnoreError(recvBuff) {
				case packet.Ping:
					select {
					case <-c.closeNotif:
					case c.sendNotif <- packet.NewPong().ToBytesIgnoreError():
					}
				case packet.Pong:
				default:
					select {
					case <-c.closeNotif:
					case c.recvNotif <- [][]byte{recvBuff}:
					}
				}
				recvBuff = nil
			} else {
				recvBuff = append(recvBuff, msg...)
			}
		}
	}()
	go func() {
		defer func() { _ = c.conn.Close() }()
		for c.active {
			select {
			case <-c.closeNotif:
				err := c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
				if err != nil {
					c.closeMutex.Lock()
					defer c.closeMutex.Unlock()
					if c.active {
						c.active = false
						_ = c.close(err)
					}
					return
				}
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			case msg := <-c.sendNotif:
				err := c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
				if err != nil {
					c.closeMutex.Lock()
					defer c.closeMutex.Unlock()
					if c.active {
						c.active = false
						_ = c.close(err)
					}
					return
				}
				err = c.conn.WriteMessage(websocket.TextMessage, append(msg, []byte("\r\n")...))
				if err != nil {
					c.closeMutex.Lock()
					defer c.closeMutex.Unlock()
					if c.active {
						c.active = false
						_ = c.close(err)
					}
					return
				}
			}
		}
	}()
	c.wsKeepAliveStart()
}

func (c *Client) wsKeepAliveStart() {
	go func() {
		kAlive := time.NewTicker(c.keepAlive)
		defer kAlive.Stop()
		for c.active {
			_ = c.Send(packet.NewPing())
			kaVal := c.keepAlive
			if kaVal > 0 {
				kAlive.Reset(kaVal)
			}
			select {
			case <-kAlive.C:
			case <-c.kaNotif:
				return
			case <-c.closeNotif:
				return
			}
		}
	}()
}

func (c *Client) GetID() string {
	if c == nil {
		return ""
	}
	if c.conn == nil {
		return "rest"
	} else {
		return "ws"
	}
}

func (c *Client) IsActive() bool {
	if c == nil {
		return false
	}
	return c.active
}

func (c *Client) Send(p *packet.Packet) error {
	if !c.IsActive() {
		return errors.New("client not active")
	}
	bts, err := p.ToBytes()
	if err != nil {
		return err
	}
	debugPrintln("Client.send - bts: " + string(bts))
	select {
	case <-c.closeNotif:
		return errors.New("client closed")
	case c.sendNotif <- bts:
	}
	return nil
}

func (c *Client) Receive() (*packet.Packet, error) {
	if !c.IsActive() {
		return nil, errors.New("client not active")
	}
	c.recvMutex.Lock()
	defer c.recvMutex.Unlock()
	if len(c.recvBuffer) < 1 {
		select {
		case <-c.closeNotif:
			return nil, errors.New("client closed")
		case pl := <-c.recvNotif:
			c.recvBuffer = append(c.recvBuffer, pl...)
		}
	}
	if len(c.recvBuffer) > 0 {
		var tr []byte
		tr, c.recvBuffer = c.recvBuffer[0], c.recvBuffer[1:]
		return packet.From(tr)
	}
	return nil, errors.New("client receive empty")
}

func (c *Client) close(err error) error {
	if c.conn != nil {
		_ = c.conn.Close()
	}
	if c.restClient != nil {
		c.restClient.CloseIdleConnections()
	}
	close(c.closeNotif)
	//close(c.recvNotif)
	//close(c.sendNotif)
	//close(c.kaNotif)
	if c.closeEvent != nil {
		c.closeEvent(c, err)
	}
	c.recvBuffer = nil
	c.sendBuffer = nil
	return err
}

func (c *Client) Close() error {
	c.closeMutex.Lock()
	defer c.closeMutex.Unlock()
	if c.active {
		c.active = false
		if c.conn != nil {
			_ = c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(c.timeout))
		}
		if c.restClient != nil {
			dRestClient := &http.Client{Timeout: c.timeout}
			defer dRestClient.CloseIdleConnections()
			req, err := http.NewRequest(http.MethodDelete, c.restTargetURL, nil)
			if err != nil {
				return c.close(err)
			}
			rsp, err := dRestClient.Do(req)
			if err != nil {
				return c.close(err)
			}
			defer func() { _, _ = io.Copy(io.Discard, rsp.Body); _ = rsp.Body.Close() }()
			return c.close(nil)
		} else {
			return c.close(nil)
		}
	}
	return nil
}

func (c *Client) SetOnClose(callback func(t Transport, e error)) {
	if c == nil {
		return
	}
	c.closeEvent = callback
}

func (c *Client) SetTimeout(to time.Duration) {
	if c == nil || to < 0 {
		return
	}
	c.timeout = to
	if c.restClient != nil {
		c.restClient.Timeout = to
	}
}

func (c *Client) GetTimeout() time.Duration {
	if c == nil {
		return 0
	}
	return c.timeout
}

func (c *Client) SetKeepAlive(ka time.Duration) {
	if c == nil {
		return
	}
	if c.IsActive() {
		c.kaMutex.Lock()
		defer c.kaMutex.Unlock()
		if ka > 0 {
			if c.keepAlive < 1 && c.conn != nil {
				c.keepAlive = ka
				c.wsKeepAliveStart()
			} else {
				c.keepAlive = ka
			}
		} else if c.conn != nil {
			select {
			case <-c.closeNotif:
			case c.kaNotif <- true:
			}
			c.keepAlive = 0
		}
	} else {
		c.keepAlive = ka
	}
}

func (c *Client) GetKeepAlive() time.Duration {
	if c == nil {
		return 0
	}
	return c.keepAlive
}

func (*Client) SetReadLimit(limit int64) {
}

func (c *Client) GetReadLimit() int64 {
	return math.MaxInt64
}
