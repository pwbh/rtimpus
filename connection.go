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
	header := parseHeader(message)
	fmt.Printf("FMT: %d, chunk stream id: %d, header length: %d\n", header.BasicHeader.Type, header.BasicHeader.StreamID, header.BasicHeader.HeaderLength)
	fmt.Printf("Timestamp: %d | Message length: %d | Message ID: %d | Message Stream ID: %d\n", header.Timestamp, header.MessageLength, header.MessageTypeId, header.MessageStreamId)
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
