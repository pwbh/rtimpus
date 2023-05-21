package main

import (
	"fmt"

	"github.com/pwbh/rtimpus"
)

func main() {
	listener, err := rtimpus.Listen("localhost:1935")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Listening for RTMP connections on %s\n", listener.Addr())
	defer listener.Close()
	rtimpus.LoopConnections(listener)

}
