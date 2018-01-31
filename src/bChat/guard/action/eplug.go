package action

import (
	"bChat/guard/hub"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// EntryPlug 插入栓，供 websoket 插入
func EntryPlug(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Errorf("upgrade request to websocket error: %s", err)
		// TODO: response error
		return
	}

	eva := hub.BuildEva(conn)
	hub.GlobalHub.Load(eva)

	// start seperate goroutine handle read and write message
	go eva.Shooting()
	go eva.Absorbing()
}
