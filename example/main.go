package main

import (
	"fmt"

	"github.com/pwbh/rtimpus"
)

func main() {
	listener, _ := rtimpus.StartRTMPListener("localhost:1935")
	fmt.Printf("Listening for RTMP connections on %v\n", listener.Addr().String())
	defer rtimpus.Close(listener)
	rtimpus.LoopConnections(listener)
}
