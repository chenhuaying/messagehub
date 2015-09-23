package main

import (
	"log"
)

type Tracker struct {
	peers   Peers
	comeIn  chan *Peer
	downOff chan *Peer
	request chan Task
}

var tracker = Tracker{
	make(Peers),
	make(chan *Peer, 256),
	make(chan *Peer, 256),
	make(chan Task, 256),
}

func (t *Tracker) run() {
	log.Println("tracker running")
	for {
		select {
		case p := <-t.comeIn:
			t.peers[p.ws.RemoteAddr().String()] = p
		case p := <-t.downOff:
			delete(t.peers, p.ws.RemoteAddr().String())
			close(p.output)
		case r := <-t.request:
			if r != nil {
				log.Println("tracker receive<<<", r)
				r.doTask(t.peers)
			}
		}
	}
}
