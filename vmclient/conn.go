package vmclient

import (
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// maxMessageSize bounds the size of a WebSocket message written by conn.Write.
	maxMessageSize = 32 * 1024
	closeTimeout   = time.Second
)

// conn adapts a WebSocket carrying subresources.PlainStreamProtocol to a net.Conn.
//
// Read and Write may be called concurrently with each other, and Close may be called at any
// time. As with gorilla/websocket, a read that hits its deadline breaks the connection; later
// reads return the same error.
type conn struct {
	ws *websocket.Conn

	readMu sync.Mutex
	reader io.Reader // current message, nil between messages

	writeMu sync.Mutex

	closeOnce sync.Once
	closeErr  error
}

var _ net.Conn = &conn{}

func newConn(ws *websocket.Conn) *conn {
	return &conn{ws: ws}
}

// Read reads stream bytes, spanning WebSocket message boundaries. It returns io.EOF once the
// server closes the stream normally.
func (c *conn) Read(p []byte) (int, error) {
	c.readMu.Lock()
	defer c.readMu.Unlock()
	for {
		if c.reader == nil {
			msgType, r, err := c.ws.NextReader()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					return 0, io.EOF
				}
				return 0, err
			}
			if msgType != websocket.BinaryMessage && msgType != websocket.TextMessage {
				continue
			}
			c.reader = r
		}
		n, err := c.reader.Read(p)
		if errors.Is(err, io.EOF) {
			c.reader = nil
			if n == 0 {
				continue
			}
			err = nil
		}
		return n, err
	}
}

// Write sends p as one or more binary messages.
func (c *conn) Write(p []byte) (int, error) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	written := 0
	for written < len(p) {
		chunk := p[written:min(len(p), written+maxMessageSize)]
		if err := c.ws.WriteMessage(websocket.BinaryMessage, chunk); err != nil {
			return written, err
		}
		written += len(chunk)
	}
	return written, nil
}

// Close sends a normal close frame and closes the connection.
func (c *conn) Close() error {
	c.closeOnce.Do(func() {
		// WriteControl may run concurrently with Write.
		_ = c.ws.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(closeTimeout))
		c.closeErr = c.ws.Close()
	})
	return c.closeErr
}

func (c *conn) LocalAddr() net.Addr  { return c.ws.LocalAddr() }
func (c *conn) RemoteAddr() net.Addr { return c.ws.RemoteAddr() }

func (c *conn) SetDeadline(t time.Time) error {
	return errors.Join(c.ws.SetReadDeadline(t), c.ws.SetWriteDeadline(t))
}

func (c *conn) SetReadDeadline(t time.Time) error  { return c.ws.SetReadDeadline(t) }
func (c *conn) SetWriteDeadline(t time.Time) error { return c.ws.SetWriteDeadline(t) }
