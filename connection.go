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
	Handshake HandshakeState
	Err       error
	conn      *net.TCPConn
}

func (c *Connection) ProcessMessage(message []byte) {
	switch c.Handshake {
	case Uninitialized:
		c.processUninitialized(message)
	case VersionSent:
		c.processVersionSent(message)
	case AckSent:
		c.processAckSent(message)
	case Done: // client and the server exchange messages.
		c.processDone(message)
	}
}

func (c *Connection) processDone(message []byte) {
	// Exchange of messages happens here.
	println("Exchanging messages?")
}

func (c *Connection) processAckSent(message []byte) {
	// C2 is received here if everything went as expected

	fmt.Println("C2:")
	fmt.Println(message[:1536])
	c.Handshake = Done
}

func (c *Connection) processVersionSent(message []byte) {

	s2 := make([]byte, 8, HANDSHAKE_PACKET_SIZE)
	clientTimestamp := message[:4]
	copy(s2, clientTimestamp)
	binary.BigEndian.PutUint32(s2[4:], uint32(time.Now().Unix()))
	hash := message[8:]
	s2 = append(s2, hash...)

	// fmt.Println(s2)

	c.conn.Write(s2)

	c.Handshake = AckSent
}

func (c *Connection) processUninitialized(message []byte) {
	if isVersionSupported(message) {
		// S0
		s0 := []byte{SUPPORTED_PROTOCOL_VERSION}
		c.conn.Write(s0)

		// S1
		s1 := make([]byte, 0, HANDSHAKE_PACKET_SIZE)
		// Start of stream timestamp 0
		s1 = append(s1, []byte{0, 0, 0, 0}...)
		// 4 bytes of 0s
		s1 = append(s1, []byte{0, 0, 0, 0}...)
		hash := utils.RandString(HANDSHAKE_PACKET_SIZE - 8)
		s1 = append(s1, hash...)

		fmt.Println("S1:")
		fmt.Println(s1)

		c.conn.Write(s1)

		c.Handshake = VersionSent
		// Sending S2 as we've received C1 already
		c.processVersionSent(message[1:])
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
