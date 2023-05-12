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
	// header := parseHeader(message)
	fmt.Printf("Message len: %d\n", len(message))
	totalByteParsed := 0

	for totalByteParsed < len(message) {
		chunk := parseChunk(message[totalByteParsed:])
		fmt.Printf("Chunk Type: %d | Chunk Stream ID: %d | Timestamp: %d | Message Length: %d | Message Type ID: %d | Message Stream ID: %d\n", chunk.header.BasicHeader.Type, chunk.header.BasicHeader.StreamID, chunk.header.Timestamp, chunk.header.MessageLength, chunk.header.MessageTypeId, chunk.header.MessageStreamId)
		fmt.Println(chunk.payload.data)

		if chunk.header.MessageTypeId == 20 {
			UnmarshalAMF0(chunk.payload.data)
		}

		totalByteParsed += int(getChunkHeaderLength(chunk.header) + chunk.header.MessageLength)
	}

	// chunk := parseChunk(message)
	// fmt.Printf("Chunk Type: %d | Chunk Stream ID: %d | Timestamp: %d | Message Length: %d | Message Type ID: %d | Message Stream ID: %d\n", chunk.header.BasicHeader.Type, chunk.header.BasicHeader.StreamID, chunk.header.Timestamp, chunk.header.MessageLength, chunk.header.MessageTypeId, chunk.header.MessageStreamId)
	// fmt.Printf("New chunk size: %d (%d)\n", chunk.payload.data, binary.BigEndian.Uint32(chunk.payload.data))
	// fmt.Println(message)
}

func (c *Connection) processAckSent(message []byte) {
	// skipping 8 bytes - 4 bytes time, 4 bytes time 2 bytes
	randomEcho := message[8:HANDSHAKE_PACKET_SIZE]

	if verifyRandomData([]byte(c.Hash), randomEcho) {
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
		binary.BigEndian.PutUint32(s1, uint32(time.Now().Unix()))
		// 4 bytes of 0s
		s1 = append(s1, []byte{0, 0, 0, 0}...)
		hash := utils.RandString(HANDSHAKE_PACKET_SIZE - 8)
		s1 = append(s1, hash...)
		c.conn.Write(s1)

		c.Phase = VersionSent
		c.Hash = hash

		c1 := message[1:]

		// S2
		s2 := make([]byte, 8, HANDSHAKE_PACKET_SIZE)
		copy(s2, c1[:4])
		// Setting time2 to when we started reading from the beginning of the handshake e.g. no time passed yet (0)
		binary.BigEndian.PutUint32(s2[4:], 0)
		c1Hash := c1[8:]
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
