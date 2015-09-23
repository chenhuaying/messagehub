package main

import (
	"encoding/json"
	"log"
	"strings"
)

//func NewBroadcast(message []byte, p *Peer) *CheckList {
//type ParseTable map[string]func(message []byte, p *Peer) Task
type ParseTable map[string]func(data map[string]interface{}, p *Peer) Task

var parseTable = ParseTable{
	//"broadcast": NewBroadcast,
	"broadcast": genBroadcastTask,
}

type SimpleProtocol struct {
	Opt  string      `json:"opt"`
	Data interface{} `json:"data"`
}

// Message string `json:"message"`
type BroadcastMessage map[string]interface{}

func genBroadcastTask(data map[string]interface{}, p *Peer) Task {
	msg := []byte(data["message"].(string))
	return NewBroadcast(msg, p)
}

//parse(message []byte, peer *Peer) CheckList
func (p *SimpleProtocol) parse(raw []byte, peer *Peer) Task {
	if err := json.Unmarshal(raw, p); err != nil {
		log.Println("SimpleProtocol parse error:", err)
		return nil
	}
	//log.Println(p.Opt, p.Data.(map[string]interface{})["message"])
	//log.Println(p.Opt, p.Data.(BroadcastMessage)["message"])
	//message := p.Data.(map[string]interface{})["message"].(string)
	message := p.Data.(map[string]interface{})
	opt := strings.ToLower(p.Opt)

	processor := parseTable[opt]
	if processor != nil {
		return processor(message, peer)
	} else {
		return nil
	}
}
