package main

import (
	"log"
)

type TaskType int

const (
	BROADCAST TaskType = iota
	REGISTER
	UPDATE
	QUERY
)

type CheckList struct {
	data interface{}
	peer *Peer
}

type Task interface {
	doTask(Peers) error
}

type Broadcast struct {
	data []byte
	peer *Peer
}

func (b *Broadcast) doTask(peers Peers) error {
	for _, peer := range peers {
		if b.peer == peer {
			log.Println("Do not send to self, Peer", b.peer.uid)
			continue
		}
		select {
		case peer.output <- b.data:
		default:
			// XXX TODO: just block sending, close is right?
			log.Println("sending block close peer")
			close(peer.output)
			delete(peers, peer.ws.RemoteAddr().String())
		}
	}
	return nil
}

func NewBroadcast(message []byte, p *Peer) Task {
	return &Broadcast{data: message, peer: p}
}
