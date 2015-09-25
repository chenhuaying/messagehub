package main

import (
	"crypto/sha1"
	"encoding/binary"
	"log"
)

const (
	BUCKET_SIZE = 511
)

type Container struct {
	peers Peers
}

type Bucket struct {
	Containers map[string]*Container
}

type BucketPool struct {
	Buckets []*Bucket
}

type Tracker struct {
	peers      Peers
	bucketPool *BucketPool
	comeIn     chan *Peer
	downOff    chan *Peer
	request    chan Task
}

var tracker = Tracker{
	make(Peers),
	NewBuckerPool(),
	make(chan *Peer, 256),
	make(chan *Peer, 256),
	make(chan Task, 256),
}

func genBucketNum(channelId string) uint32 {
	hash := sha1.Sum([]byte(channelId))
	tmp := make([]byte, len(hash))
	for idx, item := range hash {
		tmp[idx] = item
	}
	hashNum := binary.LittleEndian.Uint32(tmp)
	return hashNum % BUCKET_SIZE
}

func (t *Tracker) run() {
	log.Println("tracker running")
	for {
		select {
		case p := <-t.comeIn:
			t.peers[p.ws.RemoteAddr().String()] = p
		case p := <-t.downOff:
			delete(t.peers, p.ws.RemoteAddr().String())
			// XXX TODO: delete it from channel peer pool
		case r := <-t.request:
			if r != nil {
				log.Println("tracker receive<<<", r)
				// check registered or not
				if r.valid() {
					group := genBucketNum(r.peerCid())
					log.Println("group:", group)
					bucket := t.bucketPool.Buckets[group]
					// this Container  just a local variable
					if _, ok := bucket.Containers[r.peerCid()]; !ok {
						bucket.Containers[r.peerCid()] = new(Container)
					}
					container := bucket.Containers[r.peerCid()]
					log.Println("Cid:", r.peerCid(), container)
					if container.peers == nil {
						container.peers = make(Peers)
					}
					peers := container.peers
					log.Println("peers:", peers)
					//r.doTask(t.peers)
					r.doTask(peers)
				}
			}
		}
	}
}

func NewBuckerPool() *BucketPool {
	pool := &BucketPool{make([]*Bucket, BUCKET_SIZE)}
	for i := 0; i < BUCKET_SIZE; i++ {
		pool.Buckets[i] = new(Bucket)
		pool.Buckets[i].Containers = make(map[string]*Container)
	}
	return pool
}
