package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	request := `{"opt": "Broadcast", "data": {"message": "hello world"}}`
	msg := &SimpleProtocol{}
	peer := Peer{}
	res := msg.parse([]byte(request), &peer)
	if res == nil {
		t.Errorf("parse error, have no Broadcast protocol")
	}
}
