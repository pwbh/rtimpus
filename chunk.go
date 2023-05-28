package rtimpus

import (
	"encoding/binary"
	"fmt"

	"github.com/pwbh/rtimpus/utils"
)

type BasicHeader struct {
	Type     uint8
	StreamID uint32
	Length   uint32
}

type MessageHeader struct {
	Timestamp uint32
	Length    uint32
	TypeID    uint16
	StreamID  uint32
}

type Header struct {
	BasicHeader   *BasicHeader
	MessageHeader *MessageHeader
}

type Payload struct {
	data []byte
}

type Chunk struct {
	header  *Header
	payload *Payload
}

func (c *Chunk) Size() uint32 {
	return getChunkHeaderLength(c.header) + c.header.MessageHeader.Length
}

func parseChunk(c *Connection, message []byte) (*Chunk, error) {
	header, err := parseHeader(c, message)
	if err != nil {
		return nil, err
	}
	payload := getPayload(message, header)
	return &Chunk{
		header:  header,
		payload: payload,
	}, nil
}

func getPayload(message []byte, header *Header) *Payload {
	chunkHeaderLength := getChunkHeaderLength(header)
	max := uint32(0)
	messageLen := uint32(len(message))
	if chunkHeaderLength+header.MessageHeader.Length <= messageLen {
		max = chunkHeaderLength + header.MessageHeader.Length
	} else {
		max = messageLen
	}
	return &Payload{data: message[chunkHeaderLength:max]}
}

func parseHeader(c *Connection, message []byte) (*Header, error) {
	basicHeader := parseBasicHeader(message)

	switch basicHeader.Type {
	case 0:
		timestamp := uint32(uint(message[basicHeader.Length+2]) | uint(message[basicHeader.Length+1])<<8 | uint(message[basicHeader.Length])<<16)
		messageLength := uint32(uint(message[basicHeader.Length+5])|uint(message[basicHeader.Length+4])<<8|uint(message[basicHeader.Length+3])<<16) + 1
		messageTypeId := binary.BigEndian.Uint16([]byte{0x00, message[basicHeader.Length+6]})
		messageStreamId := binary.LittleEndian.Uint32([]byte{message[basicHeader.Length+7], message[basicHeader.Length+8], message[basicHeader.Length+9], message[basicHeader.Length+10]})

		messageHeader := &MessageHeader{Timestamp: timestamp, Length: messageLength, TypeID: messageTypeId, StreamID: messageStreamId}

		return &Header{
			BasicHeader:   basicHeader,
			MessageHeader: messageHeader,
		}, nil
	case 1:
		timestampDelta := uint32(uint(message[basicHeader.Length+2]) | uint(message[basicHeader.Length+1])<<8 | uint(message[basicHeader.Length])<<16)
		messageLength := uint32(uint(message[basicHeader.Length+5])|uint(message[basicHeader.Length+4])<<8|uint(message[basicHeader.Length+3])<<16) + 1
		messageTypeId := binary.BigEndian.Uint16([]byte{0x00, message[basicHeader.Length+6]})

		messageHeader := &MessageHeader{Timestamp: timestampDelta, Length: messageLength, TypeID: messageTypeId, StreamID: c.PrevChunk.header.BasicHeader.StreamID}

		return &Header{
			BasicHeader:   basicHeader,
			MessageHeader: messageHeader,
		}, nil
	case 2:
		timestampDelta := uint32(uint(message[basicHeader.Length+2]) | uint(message[basicHeader.Length+1])<<8 | uint(message[basicHeader.Length])<<16)

		messageHeader := &MessageHeader{Timestamp: timestampDelta, Length: c.PrevChunk.header.BasicHeader.Length, StreamID: c.PrevChunk.header.BasicHeader.StreamID}

		return &Header{
			BasicHeader:   basicHeader,
			MessageHeader: messageHeader,
		}, nil
	case 3:
		return &Header{
			BasicHeader:   basicHeader,
			MessageHeader: c.PrevChunk.header.MessageHeader,
		}, nil
	default:
		return nil, fmt.Errorf("no header available with given header type %d", basicHeader.Type)
	}
}

func parseBasicHeader(b []byte) *BasicHeader {
	basicHeader := new(BasicHeader)

	basicHeader.Type = uint8(b[0] >> 6)        // get the first two bits representing fmt
	basicHeader.StreamID = uint32(b[0] & 0x3f) // get the least significant 6 bits representing the chunk stream ID
	basicHeader.Length = 1                     // initialize header length to 1 byte

	// StreamID == 2 - reserved for low-level protocol control messages and commands
	// StreamID == - 3 - 64 represent the complete stream ID // e.g. 1 byte only

	if basicHeader.StreamID == 0 {
		basicHeader.StreamID = uint32(b[1]) + 64
		basicHeader.Length = 2
	} else if basicHeader.StreamID == 1 {
		basicHeader.StreamID = uint32(b[2])*256 + uint32(b[1]) + 64
		basicHeader.Length = 3
	}
	return basicHeader
}

func getChunkHeaderLength(header *Header) uint32 {
	switch header.BasicHeader.Type {
	case 0:
		return 11 + header.BasicHeader.Length
	case 1:
		return 7 + header.BasicHeader.Length
	case 2:
		return 3 + header.BasicHeader.Length
	case 3:
		return header.BasicHeader.Length
	default:
		return 0
	}
}

func createChunkHeader(chunkStreamId uint8, messageStreamId uint32, messageTypeID byte, payloadLength int) ([]byte, error) {
	buf := make([]byte, 12)
	buf[0] = chunkStreamId                              // Chunk Stream ID
	if err := utils.PutUint24(buf[1:], 0); err != nil { // Timestamp is ignored
		return nil, err
	}
	if err := utils.PutUint24(buf[4:], uint32(payloadLength)); err != nil { // Message length
		return nil, err
	}
	buf[7] = messageTypeID                                  // Message Type
	binary.LittleEndian.PutUint32(buf[8:], messageStreamId) // Message Stream ID
	return buf, nil
}
