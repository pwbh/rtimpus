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
}

func (c *Connection) processAckSent(message []byte) {
	// C2 is received here if everything went as expected
	fmt.Println(message)
	c.Handshake = Done
}

func (c *Connection) processVersionSent(message []byte) {
	// Send S2
	s2 := make([]byte, 0, HANDSHAKE_PACKET_SIZE)

	timestamp := time.Now().Unix()
	binary.BigEndian.PutUint32(s2, uint32(timestamp))
	clientTimestamp := message[4:8]
	s2 = append(s2, clientTimestamp...)
	hash := message[8:]
	s2 = append(s2, hash...)

	c.conn.Write(s2)
	c.Handshake = AckSent
}

func (c *Connection) processUninitialized(message []byte) {
	if isVersionSupported(message) {
		// Send S0
		c.conn.Write([]byte{SUPPORTED_PROTOCOL_VERSION})

		// Send S1
		s1 := make([]byte, 0, HANDSHAKE_PACKET_SIZE)

		timestamp := time.Now().Unix()
		binary.BigEndian.PutUint32(s1, uint32(timestamp))
		s1 = append(s1, []byte{0, 0, 0, 0}...)
		hash := []byte(utils.RandString(HANDSHAKE_PACKET_SIZE - 8))
		s1 = append(s1, []byte(hash)...)

		c.conn.Write(s1)
		c.Handshake = VersionSent
	} else {
		c.conn.Write([]byte{SUPPORTED_PROTOCOL_VERSION})
	}
}

func isVersionSupported(message []byte) bool {
	if len(message) < 1 {
		return false
	}

	return SUPPORTED_PROTOCOL_VERSION == message[0]
}
