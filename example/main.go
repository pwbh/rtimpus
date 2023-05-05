package main

import (
	"fmt"

	"github.com/pwbh/rtimpus"
)

func main() {
	listener, _ := rtimpus.StartRTMPListener("localhost:8000")
	fmt.Println("Starting a test rtimpus listener")
	defer rtimpus.Close(listener)
	rtimpus.LoopConnections(listener)
}
