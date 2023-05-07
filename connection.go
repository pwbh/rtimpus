package rtimpus

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/pwbh/rtimpus/utils"
)

const HANDSHAKE_PACKET_SIZE = int(1536)

type Connection struct {
	Phase Phase
	Err   error
	conn  *net.TCPConn
	Hash  string
}

func (c *Connection) Process(message []byte) {
	switch c.Phase {
	case Uninitialized:
		c.processUninitialized(message)
	case AckSent:
		c.processAckSent(message)
	case HandshakeDone: // client and the server exchange messages.
		c.handleChunk(message)
	}
}

func (c *Connection) handleChunk(message []byte) {
	// Exchange of messages happens here.
	fmt.Printf("Chunk arrived of size %d bytes\n", len(message))
	fmt.Printf("Chunk stream id: %d fmt: %d\n", getChunkStreamID(message), getFmt(message))

}

func getChunkStreamID(b []byte) uint32 {
	var chunkStreamID uint32
	firstByte := uint32(b[0])
	if firstByte < 64 {
		chunkStreamID = firstByte
	} else if firstByte < 128 {
		secondByte := uint32(b[1])
		chunkStreamID = (firstByte-64)<<8 + secondByte + 64
	} else {
		thirdByte := uint32(b[2])
		chunkStreamID = (firstByte-128)<<16 + uint32(b[1])<<8 + thirdByte + 64
	}
	return chunkStreamID
}

func getFmt(b []byte) uint8 {
	return uint8(b[0] >> 6)
}

func (c *Connection) processAckSent(message []byte) {
	if verifyRandomData([]byte(c.Hash), message[8:HANDSHAKE_PACKET_SIZE]) {
		c.Phase = HandshakeDone
		c.handleChunk(message[HANDSHAKE_PACKET_SIZE:])
	} else {
		c.conn.Close()
	}
}

func (c *Connection) processUninitialized(message []byte) {
	if isVersionSupported(message) {
		// S0
		s0 := []byte{SUPPORTED_PROTOCOL_VERSION}
		c.conn.Write(s0)

		// S1
		s1 := make([]byte, 4, HANDSHAKE_PACKET_SIZE)
		// Start of stream timestamp 0 also can be local server time
		binary.BigEndian.PutUint32(s1, uint32(time.Now().Unix()))
		// 4 bytes of 0s
		s1 = append(s1, []byte{0, 0, 0, 0}...)
		hash := utils.RandString(HANDSHAKE_PACKET_SIZE - 8)
		s1 = append(s1, hash...)
		c.conn.Write(s1)

		c.Phase = VersionSent
		c.Hash = hash

		// S2
		s2 := make([]byte, 8, HANDSHAKE_PACKET_SIZE)
		clientTimestamp := message[1:5]
		copy(s2, clientTimestamp)
		binary.BigEndian.PutUint32(s2[4:], uint32(time.Now().Unix()))
		c1Hash := message[9:]
		s2 = append(s2, c1Hash...)
		c.conn.Write(s2)

		c.Phase = AckSent
	} else {
		c.conn.Close()
	}

}

func isVersionSupported(message []byte) bool {
	if len(message) == 0 {
		return false
	}

	return SUPPORTED_PROTOCOL_VERSION == message[0]
}

func verifyRandomData(expected, result []byte) bool {
	if len(expected) != len(result) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != result[i] {
			return false
		}
	}

	return true
}
