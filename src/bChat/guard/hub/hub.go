package hub

import log "github.com/sirupsen/logrus"

// Hub store clients map and process message from client
type Hub struct {
	// map[uid]eva
	launchpad map[string]*Eva
	// transfer message from eva to nerv
	transport chan *Message
	load      chan *Eva
	unload    chan *Eva
}

// Message Hub 传递的消息
type Message struct {
	EvaUID string
	Msg    []byte
}

// GlobalHub 全局 HUB
var GlobalHub *Hub

// Load register client
func (h *Hub) Load(e *Eva) {
	h.load <- e
}

// Unload unregister client
func (h *Hub) Unload(e *Eva) {
	h.unload <- e
}

// Transport transfer message to hub
func (h *Hub) Transport(msg *Message) {
	h.transport <- msg
}

func init() {
	GlobalHub = &Hub{
		load:      make(chan *Eva),
		unload:    make(chan *Eva),
		transport: make(chan *Message),
		launchpad: make(map[string]*Eva),
	}

	go GlobalHub.run()
}

func (h *Hub) run() {
	for {
		select {
		case eva := <-h.load:
			h.launchpad[eva.uid] = eva
		case eva := <-h.unload:
			if _, ok := h.launchpad[eva.uid]; ok {
				delete(h.launchpad, eva.uid)
				close(eva.shoot)
			}
		case msg := <-h.transport:
			if eva, ok := h.launchpad[msg.EvaUID]; ok {
				select {
				case eva.shoot <- msg.Msg:
				default:
					log.Errorf("eva uid: %s shoot chan blocked, close and delete it", eva.uid)
					close(eva.shoot)
					delete(h.launchpad, eva.uid)
				}
			}
		}
	}
}
