package rtimpus

import (
	"encoding/binary"
	"fmt"
)

type BasicHeader struct {
	Type         uint8
	StreamID     uint32
	HeaderLength uint32
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

	headerLength := basicHeader.HeaderLength

	switch basicHeader.Type {
	case 0:
		timestamp := uint32(uint(message[headerLength+3]) | uint(message[headerLength+5])<<2 | uint(message[headerLength+1])<<16)
		messageLength := uint32(uint(message[headerLength+6]) | uint(message[headerLength+5])<<8 | uint(message[headerLength+4])<<16)
		messageTypeId := binary.BigEndian.Uint16([]byte{0x00, message[headerLength+7]})
		messageStreamId := binary.BigEndian.Uint32([]byte{message[headerLength+8], message[headerLength+9], message[headerLength+10], message[headerLength+11]})

		return &Header{
			BasicHeader:     basicHeader,
			Timestamp:       timestamp,
			MessageLength:   messageLength,
			MessageTypeId:   messageTypeId,
			MessageStreamId: messageStreamId,
		}
	case 1:
		timestampDelta := uint32(uint(message[headerLength+3]) | uint(message[headerLength+5])<<2 | uint(message[headerLength+1])<<16)
		messageLength := uint32(uint(message[headerLength+6]) | uint(message[headerLength+5])<<8 | uint(message[headerLength+4])<<16)
		messageTypeId := binary.BigEndian.Uint16([]byte{0x00, message[headerLength+7]})

		return &Header{
			BasicHeader:   basicHeader,
			Timestamp:     timestampDelta,
			MessageLength: messageLength,
			MessageTypeId: messageTypeId,
		}
	case 2:
		timestampDelta := uint32(uint(message[headerLength+3]) | uint(message[headerLength+5])<<2 | uint(message[headerLength+1])<<16)

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
	basicHeader.HeaderLength = 1               // initialize header length to 1 byte

	if basicHeader.StreamID == 0 {
		basicHeader.StreamID = uint32(b[1]) + 64
		basicHeader.HeaderLength = 2
	} else if basicHeader.StreamID == 1 {
		basicHeader.StreamID = uint32(b[2])*256 + uint32(b[1]) + 64
		basicHeader.HeaderLength = 3
	}

	return basicHeader
}

func getChunkMessageHeaderSizeByType(t uint8) (uint8, error) {
	switch t {
	case 0:
		return 11, nil
	case 1:
		return 7, nil
	case 2:
		return 3, nil
	case 3:
		return 0, nil

	default:
		return 0, fmt.Errorf("unsupported fmt has been given")
	}
}
