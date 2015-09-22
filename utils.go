package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

var inc int32 = 0
var mtx = sync.Mutex{}

func getHostHash() []byte {
	var macAddr string
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Poor soul, error: " + err.Error())
	}
	for _, inter := range interfaces {
		mac := inter.HardwareAddr
		macAddr := mac.String()
		if macAddr != "" {
			break
		}
	}
	hostHash := md5.New()
	io.WriteString(hostHash, macAddr)
	return hostHash.Sum(nil)
}

func getTimeStamp() int64 {
	now := time.Now()
	return now.Unix()
}

func getPid() int {
	return os.Getpid()
}

func getInc() int32 {
	var res int32 = 0
	{
		mtx.Lock()
		defer mtx.Unlock()
		inc++
		res = inc
	}
	runtime.Gosched()
	return res
}

func isLittleEndian() bool {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	fmt.Println(*pb)
	return (b == 0x04)
}

// 0~3 time stamp
// 4~6 host name hash
// 7~8 process id
// 9~11 accumulator
func getPeerId() []byte {
	peerid := make([]byte, 0, 12)

	timeStampBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(timeStampBytes, uint32(getTimeStamp()))
	peerid = append(peerid, timeStampBytes[0:4]...)

	hostHash := getHostHash()
	peerid = append(peerid, hostHash[0:3]...)

	pidBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(pidBytes, uint16(getPid()))
	peerid = append(peerid, pidBytes...)

	incBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(incBytes, uint32(getInc()))
	peerid = append(peerid, incBytes[0:3]...)

	return peerid
}
