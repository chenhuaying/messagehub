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
	"broadcast":  genBroadcastTask,
	"register":   genRegisterTask,
	"unregister": genUnregisterTask,
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

func genRegisterTask(data map[string]interface{}, p *Peer) Task {
	channelId := []byte(data["channelid"].(string))
	return NewRegister(channelId, p)
}

func genUnregisterTask(data map[string]interface{}, p *Peer) Task {
	channelId := data["channelid"].(string)
	return NewUnregister(channelId, p)
}

//parse(message []byte, peer *Peer) CheckList
func (p *SimpleProtocol) parse(raw []byte, peer *Peer) Task {
	// clear
	p.Data = nil
	p.Opt = ""

	if err := json.Unmarshal(raw, p); err != nil {
		log.Println("SimpleProtocol parse error:", err)
		return nil
	}
	//log.Println(p.Opt, p.Data.(map[string]interface{})["message"])
	//log.Println(p.Opt, p.Data.(BroadcastMessage)["message"])
	//message := p.Data.(map[string]interface{})["message"].(string)
	if p.Opt != "" && p.Data != nil {
		if message, ok := p.Data.(map[string]interface{}); ok {
			opt := strings.ToLower(p.Opt)
			log.Println("Parse Protocol OK:", opt, message)
			processor := parseTable[opt]
			if processor != nil {
				return processor(message, peer)
			} else {
				// no implemention
				return nil
			}
		} else {
			return nil
		}
	}

	return nil
}
