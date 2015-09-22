package main

import (
	"log"
)

type Tracker struct {
	peers   map[string]*Peer
	comeIn  chan *Peer
	downOff chan *Peer
	request chan []byte
}

var tracker = Tracker{
	make(map[string]*Peer),
	make(chan *Peer, 256),
	make(chan *Peer, 256),
	make(chan []byte, 256),
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
			log.Println("tracker receive<<<", r)
			for _, peer := range t.peers {
				select {
				case peer.output <- r:
				default:
					log.Println("sending block close peer")
					close(peer.output)
					delete(t.peers, peer.ws.RemoteAddr().String())
				}
			}
		}
	}
}
