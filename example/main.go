package main

import (
	"fmt"

	"github.com/pwbh/rtimpus"
)

func main() {
	listener, err := rtimpus.StartRTMPListener("localhost:1935")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Listening for RTMP connections on %s\n", listener.Addr())
	defer rtimpus.Close(listener)
	rtimpus.LoopConnections(listener)
}
