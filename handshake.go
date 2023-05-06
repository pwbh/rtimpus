package rtimpus

import (
	"fmt"
)

type Phase byte

const (
	Uninitialized Phase = iota
	VersionSent
	AckSent
	HandshakeDone
)

type C0S0 struct {
	Version int
}

func Trying() {
	fmt.Println("Test")
}

//	Version defined by this specification is 3 (This is the only version this protocol will
//	support 0-2 are depercated values sused by earleir proprietary products,
//	4-31 are reserved for future implementations, 32-35 are not allowed.)

//	TODO: Handeshake
// 		C0 - RTMP version requestby the client from the server.
// 			S0 - Selected version of the RTMP version from the server.
// 				If server doesn't recognize the version sent by the client in C0
// 				server should respond to the client with 3. Client may choose
//				to degrade to version 3, or to abandon the handshake.
// 		C1 - q
