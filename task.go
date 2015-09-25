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
	peerUid() string
	peerCid() string
	valid() bool
}

//===================
// Do broadcast task
//===================
type Broadcast struct {
	data []byte
	peer *Peer
}

func (p *Peer) peerUid() string {
	return p.uid
}

func (b *Broadcast) doTask(peers Peers) error {
	for _, peer := range peers {
		if b.peer == peer {
			log.Println("Do not send to self, Peer", b.peer.uid)
			continue
		}
		select {
		case peer.output <- b.data:
		case <-peer.done:
		default:
			// XXX TODO: just block sending, close is right?
			log.Println("sending block close peer")
			close(peer.output)
			delete(peers, peer.ws.RemoteAddr().String())
		}
	}
	return nil
}

func (b *Broadcast) peerUid() string {
	return b.peer.uid
}

func (b *Broadcast) peerCid() string {
	return b.peer.cid
}

func (b *Broadcast) valid() bool {
	return isRegistered(b.peer)
}

func NewBroadcast(message []byte, p *Peer) Task {
	return &Broadcast{data: message, peer: p}
}

//===================
// Do register task
//===================
type Register struct {
	channelId string
	peer      *Peer
}

func (r *Register) doTask(peers Peers) error {
	if _, found := peers[r.peer.uid]; !found {
		log.Println("add peer:", r.peer.uid)
		// add peer to channel peer map
		peers[r.peer.uid] = r.peer
	}
	log.Println("doTask peers:", peers)
	// set peer channleId
	r.peer.cid = r.channelId

	return nil
}

func (r *Register) peerUid() string {
	return r.peer.uid
}

func (r *Register) peerCid() string {
	return r.channelId
}

func (r *Register) valid() bool {
	// XXX TODO unregister-first: if peer's channelid != register's channel id,
	//then unregister old channel and register new
	return true
}

func NewRegister(id []byte, p *Peer) Task {
	return &Register{channelId: string(id), peer: p}
}

//===================
// Do unregister task
//===================
type Unregister struct {
	channelId string
	peer      *Peer
}

func (u *Unregister) doTask(peers Peers) error {
	delete(peers, u.peer.uid)
	// clear peer channelId
	u.peer.cid = ""
	return nil
}

func (u *Unregister) peerUid() string {
	return u.peer.uid
}

func (u *Unregister) peerCid() string {
	return u.channelId
}

func (u *Unregister) valid() bool {
	// it can be true when unregister-first implemented
	if u.peer.cid != u.channelId {
		return false
	}
	return true
}

func NewUnregister(id string, p *Peer) Task {
	return &Unregister{channelId: id, peer: p}
}

// check peer is registered or not
func isRegistered(p *Peer) bool {
	if p.cid == "" {
		return false
	} else {
		return true
	}
}
