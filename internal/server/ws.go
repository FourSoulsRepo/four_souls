package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// sendQueue is how many messages wait for a slow client before it is
// disconnected; it reconnects and resyncs (ADR 006).
const sendQueue = 64

// writeWait bounds one write to a client.
const writeWait = 10 * time.Second

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	// LAN trust (ADR 006): any page may connect; a browser client is
	// served from elsewhere later.
	CheckOrigin: func(*http.Request) bool { return true },
}

// Handler serves the room over WebSocket at any path it is mounted on.
func Handler(r *Room) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ws, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			return // the upgrader already answered with an error
		}
		serveWS(r, ws)
	})
}

// wsConn is a Conn over one WebSocket: a queue and a writer goroutine.
type wsConn struct {
	ws    *websocket.Conn
	out   chan []byte
	once  sync.Once
	close chan struct{}
}

func serveWS(r *Room, ws *websocket.Conn) {
	c := &wsConn{ws: ws, out: make(chan []byte, sendQueue), close: make(chan struct{})}
	go c.write()
	client := r.Connect(c)
	defer client.Leave()
	defer c.Close()
	ws.SetReadLimit(protocol.MaxClientMessage)
	for {
		typ, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}
		if typ == websocket.TextMessage {
			client.Receive(msg)
		}
	}
}

// Send queues a message; a full queue means the client is too slow.
func (c *wsConn) Send(msg []byte) bool {
	select {
	case <-c.close:
		return false
	default:
	}
	select {
	case c.out <- msg:
		return true
	default:
		return false
	}
}

// Close ends the connection once.
func (c *wsConn) Close() {
	c.once.Do(func() {
		close(c.close)
		_ = c.ws.Close() //nolint:errcheck // closing; nothing to do on failure
	})
}

func (c *wsConn) write() {
	for {
		select {
		case <-c.close:
			return
		case msg := <-c.out:
			if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.Close()
				return
			}
			if err := c.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				c.Close()
				return
			}
		}
	}
}
