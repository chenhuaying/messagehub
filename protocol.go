package main

type Parser interface {
	parse(message []byte, peer *Peer) Task
}
