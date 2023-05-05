package main

import (
	"fmt"

	"github.com/pwbh/rtimpus"
)

func main() {
	listener, _ := rtimpus.StartRTMPListener("localhost:3333")
	fmt.Println("Starting a test RTIMPUS listener")
	defer rtimpus.Close(listener)
	rtimpus.LoopConnections(listener)
}
