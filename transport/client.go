package transport

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/gorilla/websocket"
	"golang.local/gc-c-com/packet"
	"net/http"
	"net/url"
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
	closeEvent    CloseCallback
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
		return false
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusOK && resp.ContentLength > 0 && strings.EqualFold(resp.Header.Get("Content-Type"), "text/plain") {
		buff := make([]byte, resp.ContentLength)
		ln, err := resp.Body.Read(buff)
		if err != nil {
			return false
		}
		c.restTargetURL = restURL + "?s=" + url.QueryEscape(strings.Trim(string(buff[:ln]), "\r\n"))
		if c.keepAlive < 1 {
			c.keepAlive = time.Second
		}
		go func() {
			for c.active {
				bts, ok := <-c.sendNotif
				c.appendToSendBuffer(bts)
				if !ok {
					return
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
					c.sendNotif <- packet.NewPong().ToBytesIgnoreError()
				}
				kAlive.Reset(c.keepAlive)
				select {
				case <-kAlive.C:
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

func (c *Client) appendToSendBuffer(bts []byte) {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()
	c.sendBuffer = append(c.sendBuffer, bts)
}

func (c *Client) handlerProcessor() (failed bool, hasPing bool) {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()
	var resp *http.Response
	var err error
	defer func() {
		_ = resp.Body.Close()
		c.sendBuffer = nil
	}()
	if len(c.sendBuffer) < 1 {
		resp, err = c.restClient.Get(c.restTargetURL)
	} else {
		buff := bytes.NewBuffer(make([]byte, c.getSendSize()))
		for _, bts := range c.sendBuffer {
			_, err = buff.Write(bts)
			if err != nil {
				return true, false
			}
			_, err = buff.Write([]byte("\r\n"))
			if err != nil {
				return true, false
			}
		}
		resp, err = c.restClient.Post(c.restTargetURL, "text/plain", buff)
	}
	if err != nil {
		return true, false
	}
	if resp.StatusCode == http.StatusOK && resp.ContentLength > 0 && strings.EqualFold(resp.Header.Get("Content-Type"), "text/plain") {
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
		c.recvNotif <- rIn
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
					c.sendNotif <- packet.NewPong().ToBytesIgnoreError()
				case packet.Pong:
				default:
					c.recvNotif <- [][]byte{msg}
				}
			} else if len(recvBuff) > 0 && msg[len(msg)-1] == '\n' {
				recvBuff = append(recvBuff, msg...)
				switch packet.GetCommandIgnoreError(recvBuff) {
				case packet.Ping:
					c.sendNotif <- packet.NewPong().ToBytesIgnoreError()
				case packet.Pong:
				default:
					c.recvNotif <- [][]byte{recvBuff}
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
			msg, ok := <-c.sendNotif
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
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
	}()
	c.wsKeepAliveStart()
}

func (c *Client) wsKeepAliveStart() {
	go func() {
		kAlive := time.NewTicker(c.keepAlive)
		defer kAlive.Stop()
		for c.active {
			_ = c.Send(packet.NewPing())
			kAlive.Reset(c.keepAlive)
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
	c.sendNotif <- bts
	return nil
}

func (c *Client) Receive() (*packet.Packet, error) {
	if !c.IsActive() {
		return nil, errors.New("client not active")
	}
	c.recvMutex.Lock()
	defer c.recvMutex.Unlock()
	if len(c.recvBuffer) < 1 {
		tOut := time.NewTimer(c.timeout)
		defer tOut.Stop()
		select {
		case pl := <-c.recvNotif:
			c.recvBuffer = append(c.recvBuffer, pl...)
			break
		case <-tOut.C:
			c.closeMutex.Lock()
			defer c.closeMutex.Unlock()
			if c.active {
				c.active = false
				_ = c.close(errors.New("client receive timeout"))
			}
			return nil, errors.New("client receive timeout")
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
	close(c.closeNotif)
	close(c.recvNotif)
	close(c.sendNotif)
	close(c.kaNotif)
	if c.closeEvent != nil {
		c.closeEvent(c, err)
	}
	c.recvBuffer = nil
	c.sendBuffer = nil
	return err
}

func (c *Client) Close() error {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()
	c.recvMutex.Lock()
	defer c.recvMutex.Unlock()
	c.closeMutex.Lock()
	defer c.closeMutex.Unlock()
	if c.active {
		c.active = false
		if c.conn != nil {
			_ = c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(c.timeout))
		}
		return c.close(nil)
	}
	return nil
}

func (c *Client) SetOnClose(callback CloseCallback) {
	if c == nil {
		return
	}
	c.closeEvent = callback
}

func (c *Client) SetTimeout(to time.Duration) {
	if c == nil {
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
			c.kaNotif <- true
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
