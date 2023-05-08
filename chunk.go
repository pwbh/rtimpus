package rtimpus

import (
	"encoding/binary"
)

type BasicHeader struct {
	Type     uint8
	StreamID uint32
	Length   uint32
}

type Header struct {
	BasicHeader     *BasicHeader
	Timestamp       uint32
	MessageLength   uint32
	MessageTypeId   uint16
	MessageStreamId uint32
}

func parseHeader(message []byte) *Header {
	basicHeader := parseBasicHeader(message)

	switch basicHeader.Type {
	case 0:
		timestamp := uint32(uint(message[basicHeader.Length+3]) | uint(message[basicHeader.Length+5])<<2 | uint(message[basicHeader.Length+1])<<16)
		messageLength := uint32(uint(message[basicHeader.Length+6]) | uint(message[basicHeader.Length+5])<<8 | uint(message[basicHeader.Length+4])<<16)
		messageTypeId := binary.BigEndian.Uint16([]byte{0x00, message[basicHeader.Length+7]})
		messageStreamId := binary.BigEndian.Uint32([]byte{message[basicHeader.Length+8], message[basicHeader.Length+9], message[basicHeader.Length+10], message[basicHeader.Length+11]})

		return &Header{
			BasicHeader:     basicHeader,
			Timestamp:       timestamp,
			MessageLength:   messageLength,
			MessageTypeId:   messageTypeId,
			MessageStreamId: messageStreamId,
		}
	case 1:
		timestampDelta := uint32(uint(message[basicHeader.Length+3]) | uint(message[basicHeader.Length+5])<<2 | uint(message[basicHeader.Length+1])<<16)
		messageLength := uint32(uint(message[basicHeader.Length+6]) | uint(message[basicHeader.Length+5])<<8 | uint(message[basicHeader.Length+4])<<16)
		messageTypeId := binary.BigEndian.Uint16([]byte{0x00, message[basicHeader.Length+7]})

		return &Header{
			BasicHeader:   basicHeader,
			Timestamp:     timestampDelta,
			MessageLength: messageLength,
			MessageTypeId: messageTypeId,
		}
	case 2:
		timestampDelta := uint32(uint(message[basicHeader.Length+3]) | uint(message[basicHeader.Length+5])<<2 | uint(message[basicHeader.Length+1])<<16)

		return &Header{
			BasicHeader: basicHeader,
			Timestamp:   timestampDelta,
		}

	default:
		return &Header{BasicHeader: basicHeader}
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
