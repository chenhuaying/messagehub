package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Peer struct {
	uid    string
	cid    string
	ws     *websocket.Conn
	output chan []byte
	parser Parser
}

type Peers map[string]*Peer

// write writes a message with the given message type and payload.
func (p *Peer) write(mt int, payload []byte) error {
	p.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return p.ws.WriteMessage(mt, payload)
}

func (p *Peer) readRoutine() {
	defer func() {
		log.Println("reading routine will stop")
		tracker.downOff <- p
		p.ws.Close()
	}()

	p.ws.SetReadLimit(512)
	p.ws.SetReadDeadline(time.Now().Add(pongWait))
	p.ws.SetPongHandler(func(string) error {
		log.Println("received ping message")
		p.ws.SetReadDeadline(time.Now().Add(pongWait))
		p.ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := p.ws.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
			log.Println("WebSocket write pong message failed!=", err)
			return nil
		}
		log.Println("Send Pong message")
		return nil
	})

	for {
		log.Println("waiting message .....")
		msgType, bytes, err := p.ws.ReadMessage()
		if err != nil {
			log.Println("reading message error: ", err)
			break
		}
		log.Println(msgType, string(bytes))
		task := p.parser.parse(bytes, p)
		// XXX TODO: add pretask
		tracker.request <- task
		// XXX TODO: add posttask
	}
}

func (p *Peer) writeRoutine() {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("writing routine will stop")
		pingTicker.Stop()
		p.ws.Close()
	}()

	for {
		select {
		case message, ok := <-p.output:
			if !ok {
				// output channel error, wrong closed
				p.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := p.write(websocket.TextMessage, message); err != nil {
				log.Println("write failed, error:", err)
				return
			}
		case <-pingTicker.C:
			log.Println("write ping message")
			p.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("send ping error:", err)
				return
			}
		}
	}
}

func handlerws(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	log.Println(r.RemoteAddr, r.Host)

	peer := &Peer{getHexPeerId(), "", ws, make(chan []byte, 256), &SimpleProtocol{}}
	tracker.comeIn <- peer

	go peer.writeRoutine()
	peer.readRoutine()
}
