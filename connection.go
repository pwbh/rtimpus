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
	Phase              Phase
	Err                error
	Conn               *net.TCPConn
	Hash               string
	ClientMaxChunkSize uint32
	ServerMaxChunkSize uint32
	ServerWindowSize   uint32
	totalBytesReceived uint32
	PrevChunk          *Chunk
	AccChunk           *Chunk
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

func (c *Connection) Write(b []byte) (int, error) {
	return c.Conn.Write(b)
}

func (c *Connection) handleChunk(message []byte) {
	fmt.Printf("Message len: %d\n", len(message))
	fmt.Println(message)

	currentTotalBytesReceived := uint32(0)

	for int(currentTotalBytesReceived) < len(message) {
		chunk, err := parseChunk(c, message[currentTotalBytesReceived:])

		if chunk.header.BasicHeader.Type < 3 {
			c.PrevChunk = chunk
		}

		fmt.Println(chunk.header.MessageHeader.Length, uint32(len(chunk.payload.data)))
		fmt.Println(chunk.payload.data)

		if chunk.header.MessageHeader.Length > uint32(len(chunk.payload.data)) && c.AccChunk == nil {
			fmt.Println("bigger")
			c.AccChunk = chunk
			break
		} else if c.AccChunk != nil {
			c.AccChunk.payload.data = append(c.AccChunk.payload.data, chunk.payload.data...)
			if chunk.header.MessageHeader.Length != uint32(len(c.AccChunk.payload.data)+1) {
				break
			}
			chunk = c.AccChunk
		}

		if err != nil {
			fmt.Printf("fatal error in recieving data from client: %v", err)
			return
		}
		fmt.Printf("Chunk Type: %d | Chunk Stream ID: %d | Timestamp: %d | Message Length: %d | Message Type ID: %d | Message Stream ID: %d\n", chunk.header.BasicHeader.Type, chunk.header.BasicHeader.StreamID, chunk.header.MessageHeader.Timestamp, chunk.header.MessageHeader.Length, chunk.header.MessageHeader.TypeID, chunk.header.MessageHeader.StreamID)
		switch chunk.header.MessageHeader.TypeID {
		case 1:
			c.ClientMaxChunkSize = binary.BigEndian.Uint32(chunk.payload.data)
		case 3:
			fmt.Printf("client: %d bytes acknowledged\n", binary.BigEndian.Uint32(chunk.payload.data))
		case 20, 18: // Message Type ID 20, 18 is Command Message
			command, err := UnmarshalCommand(chunk)
			if err != nil {
				fmt.Printf("error when unmarshaling command: %v\n", err)
				return
			}
			fmt.Println(command)
			c.handleCommand(command, chunk)
		default:
			fmt.Printf("message ID: %d is not handled yet\n", chunk.header.MessageHeader.TypeID)
			return
		}

		currentTotalBytesReceived += chunk.Size()

		// if err := c.checkAcknowledgement(chunk); err != nil {
		// 	fmt.Printf("ack failed: %v\n", err)
		// }
	}

	c.totalBytesReceived += currentTotalBytesReceived

}

// func (c *Connection) checkAcknowledgement(chunk *Chunk) error {
//
// 	if c.BytesRecievedNoAck >= c.ClientMaxChunkSize {
// 		diff := c.BytesRecievedNoAck - c.ClientMaxChunkSize
// 		sequenceNumber := c.BytesRecievedNoAck - diff
// 		err := sendAcknowledgement(c, sequenceNumber)
// 		if err != nil {
// 			return err
// 		}
// 		c.BytesRecievedNoAck = diff
// 	}
// 	return nil
// }

func (c *Connection) handleCommand(command interface{}, chunk *Chunk) {
	switch command.(type) {
	case *Connect:
		if err := sendWindowAcknowledgementSize(c, 4096); err != nil {
			fmt.Printf("error on sendConnectResult: %v\n", err)
		}
		if err := sendSetPeerBandwith(c, 8192, 0); err != nil {
			fmt.Printf("error on sendConnectResult: %v\n", err)
		}
		if err := sendConnectResult(c); err != nil {
			fmt.Printf("error on sendConnectResult: %v\n", err)
		}
	default:
		fmt.Printf("unrecognized command received, %v\n", command)
	}
}

func (c *Connection) processAckSent(message []byte) {
	// skipping 8 bytes - 4 bytes time, 4 bytes time 2 bytes
	randomEcho := message[8:HANDSHAKE_PACKET_SIZE]

	if verifyRandomData([]byte(c.Hash), randomEcho) {
		c.Phase = HandshakeDone
		c.handleChunk(message[HANDSHAKE_PACKET_SIZE:])
	} else {
		c.Conn.Close()
	}
}

func (c *Connection) processUninitialized(message []byte) {
	if isVersionSupported(message) {
		// S0
		s0 := []byte{SUPPORTED_PROTOCOL_VERSION}
		c.Conn.Write(s0)

		// S1
		s1 := make([]byte, 4, HANDSHAKE_PACKET_SIZE)
		binary.BigEndian.PutUint32(s1, uint32(time.Now().Unix()))
		// 4 bytes of 0s
		s1 = append(s1, []byte{0, 0, 0, 0}...)
		hash := utils.RandString(HANDSHAKE_PACKET_SIZE - 8)
		s1 = append(s1, hash...)
		c.Conn.Write(s1)

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
		c.Conn.Write(s2)

		c.Phase = AckSent
	} else {
		c.Conn.Close()
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
