package server

import (
	"encoding/binary"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// sendQueue is how many messages wait for a slow client before it is
// disconnected; it reconnects and resyncs (ADR 006).
const sendQueue = 256

// writeWait bounds one write to a client.
const writeWait = 10 * time.Second

// batchBytes is how many bytes of queued messages one write sends at
// most. Each message is still its own frame; the writes are fewer (18.17).
const batchBytes = 64 << 10

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	// LAN trust (ADR 006): any page may connect; a browser client is
	// served from elsewhere later.
	CheckOrigin: func(*http.Request) bool { return true },
}

// Handler serves the hub over WebSocket at any path it is mounted on.
func Handler(h *Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ws, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			return // the upgrader already answered with an error
		}
		serveWS(h, ws)
	})
}

// wsConn is a Conn over one WebSocket: a queue and a writer goroutine.
// The writer puts frames on the socket itself, in batches; pongs take the
// same lock so the two never interleave (18.17).
type wsConn struct {
	ws    *websocket.Conn
	conn  net.Conn // the socket under ws
	out   chan []byte
	once  sync.Once
	close chan struct{}
	local bool // the peer is on this computer
	wmu   sync.Mutex
}

func serveWS(h *Hub, ws *websocket.Conn) {
	c := &wsConn{ws: ws, conn: ws.UnderlyingConn(), out: make(chan []byte, sendQueue), close: make(chan struct{}), local: loopback(ws.RemoteAddr())}
	ws.SetPingHandler(c.pong)
	go c.write()
	client := h.Connect(c)
	defer client.Leave()
	defer c.Close()
	ws.SetReadLimit(protocol.MaxClientMessage)
	for {
		typ, msg, err := ws.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("server: read from a client ended: %v", err) // 18.17: each close has a cause
			}
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
		log.Printf("server: a client's send queue is full (%d messages); it is cut off", sendQueue)
		return false
	}
}

// Local reports whether the peer is on this computer.
func (c *wsConn) Local() bool { return c.local }

// loopback reports whether a is a loopback address.
func loopback(a net.Addr) bool {
	t, ok := a.(*net.TCPAddr)
	return ok && t.IP.IsLoopback()
}

// Close ends the connection once.
func (c *wsConn) Close() {
	c.once.Do(func() {
		close(c.close)
		_ = c.ws.Close() //nolint:errcheck // closing; nothing to do on failure
	})
}

// pong answers a ping. It takes the writer's lock, since both write to
// the socket.
func (c *wsConn) pong(data string) error {
	c.wmu.Lock()
	defer c.wmu.Unlock()
	return c.ws.WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(writeWait))
}

// write sends the queued messages. It takes every message already queued,
// up to batchBytes, and sends them in one write.
func (c *wsConn) write() {
	var batch []byte
	for {
		select {
		case <-c.close:
			return
		case msg := <-c.out:
			batch = appendFrame(batch[:0], msg)
			batch = c.takeQueued(batch)
			if err := c.flush(batch); err != nil {
				c.Close()
				return
			}
		}
	}
}

// takeQueued appends the messages waiting in the queue to batch.
func (c *wsConn) takeQueued(batch []byte) []byte {
	for len(batch) < batchBytes {
		select {
		case msg := <-c.out:
			batch = appendFrame(batch, msg)
		default:
			return batch
		}
	}
	return batch
}

// flush writes the frames to the socket.
func (c *wsConn) flush(batch []byte) error {
	c.wmu.Lock()
	defer c.wmu.Unlock()
	if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	_, err := c.conn.Write(batch)
	if err != nil {
		log.Printf("server: write to a client failed: %v", err)
	}
	return err
}

// appendFrame appends msg as one final text frame from the server, not
// masked (RFC 6455, section 5.2).
func appendFrame(b, msg []byte) []byte {
	b = append(b, 0x80|byte(websocket.TextMessage))
	switch n := len(msg); {
	case n < 126:
		b = append(b, byte(n))
	case n <= 0xffff:
		b = append(b, 126)
		b = binary.BigEndian.AppendUint16(b, uint16(n))
	default:
		b = append(b, 127)
		b = binary.BigEndian.AppendUint64(b, uint64(n))
	}
	return append(b, msg...)
}
