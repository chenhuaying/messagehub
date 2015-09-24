package main

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	request := `{"OPT": "Broadcast", "data": {"message": "hello world"}}`
	msg := &SimpleProtocol{}
	peer := Peer{}
	res := msg.parse([]byte(request), &peer)
	if res == nil {
		t.Errorf("parse error, have no Broadcast protocol")
	}

	request = `[123,456, 789]`
	res = msg.parse([]byte(request), &peer)
	if res != nil {
		t.Errorf("parse error, have no Broadcast protocol")
	}

	request = `{"Opt":"Broadcast", "data": ["come into channel 1", "hello protocol"]}`
	res = msg.parse([]byte(request), &peer)
	if res == nil {
		fmt.Println("parse error, Broadcast protocol message error")
	} else {
		if bcmsg, ok := res.(*Broadcast); ok {
			if string(bcmsg.data) != "come into channel 1" {
				t.Errorf("parse error")
			}
		}
	}

	request = `{"OPT": "Broadcast", "data": {"message": "hello here"}}`
	res = msg.parse([]byte(request), &peer)
	if res == nil {
		t.Errorf("parse error, have no Broadcast protocol")
	} else {
		if bcmsg, ok := res.(*Broadcast); ok {
			if string(bcmsg.data) != "hello here" {
				t.Errorf("parse error")
			}
		}
	}
}
