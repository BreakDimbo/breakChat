package hub

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// TODO: 使用状态机，控制客户端状态
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Eva is a middleman between the websocket connection and the hub.
type Eva struct {
	uid  string
	conn *websocket.Conn
	// send message out
	shoot chan []byte
}

// BuildEva new a client
func BuildEva(conn *websocket.Conn) *Eva {
	uid := uuid.NewV4().String()
	return &Eva{
		uid:   uid,
		conn:  conn,
		shoot: make(chan []byte, 256),
	}
}

// Shooting receive message from hub and send it to websocket connection
func (e *Eva) Shooting() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		e.conn.Close()
	}()

	for {
		select {
		case message, ok := <-e.shoot:
			e.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// the hub closed channel
				e.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := e.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// TODO: ?
			// Add queued chat messages to the current websocket message.
			n := len(e.shoot)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-e.shoot)
			}

			if err := w.Close(); err != nil {
				log.Errorf("close websocket writer error: %s", err)
				return
			}

		case <-ticker.C:
			e.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := e.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Absorbing receive message from websocket connction and send it to hub
func (e *Eva) Absorbing() {
	defer func() {
		GlobalHub.Unload(e)
		e.conn.Close()
	}()
	e.conn.SetReadLimit(maxMessageSize)
	e.conn.SetReadDeadline(time.Now().Add(pongWait))
	// TODO: understand the pong
	e.conn.SetPongHandler(func(string) error {
		e.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := e.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("error: %v", err)
			}
			break
		}

		// TODO: 定义消息协议，使用 json 传输消息。
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		m := &Message{
			EvaUID: e.uid,
			Msg:    message,
		}
		// TODO: transport to nerv by gRPC
		GlobalHub.Transport(m)
	}
}
