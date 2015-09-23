package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGetPid(t *testing.T) {
	fmt.Println(getPid())
}

func TestGetTimeStamp(t *testing.T) {
	fmt.Println(getTimeStamp())
}

//func TestGetInc(t *testing.T) {
//	go func() {
//		for i := 0; i < 20; i++ {
//			fmt.Println("1>>", getInc())
//		}
//	}()
//
//	for i := 0; i < 20; i++ {
//		fmt.Println("2>>", getInc())
//	}
//}

func TestGetPeerId(t *testing.T) {
	for i := 0; i < 10; i++ {
		id := getPeerId()
		fmt.Println(id, hex.EncodeToString(id))
		length := len(id)
		if length != 12 {
			t.Errorf("length not valide value[%d]", length)
		}
	}
}

func TestGetHexPeerId(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(getHexPeerId())
	}
}
