package rtimpus

import (
	"encoding/binary"
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
	}
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

	} else {
		c.conn.Write([]byte{SUPPORTED_PROTOCOL_VERSION})
	}
}

func isVersionSupported(message []byte) bool {
	if len(message) > 1 {
		return false
	}

	return SUPPORTED_PROTOCOL_VERSION == message[0]
}
