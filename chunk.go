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

type Payload struct {
	data []byte
}

type Chunk struct {
	header  *Header
	payload *Payload
}

func parseChunk(message []byte) *Chunk {
	header := parseHeader(message)

	chunkHeaderLength := getChunkHeaderLength(header)

	return &Chunk{
		header:  header,
		payload: &Payload{data: message[chunkHeaderLength : chunkHeaderLength+header.MessageLength]},
	}
}

func parseHeader(message []byte) *Header {
	basicHeader := parseBasicHeader(message)

	switch basicHeader.Type {
	case 0:
		timestamp := uint32(uint(message[basicHeader.Length+2]) | uint(message[basicHeader.Length+1])<<8 | uint(message[basicHeader.Length])<<16)
		messageLength := uint32(uint(message[basicHeader.Length+5]) | uint(message[basicHeader.Length+4])<<8 | uint(message[basicHeader.Length+3])<<16)
		messageTypeId := binary.BigEndian.Uint16([]byte{0x00, message[basicHeader.Length+6]})
		messageStreamId := binary.BigEndian.Uint32([]byte{message[basicHeader.Length+7], message[basicHeader.Length+8], message[basicHeader.Length+9], message[basicHeader.Length+10]})

		return &Header{
			BasicHeader:     basicHeader,
			Timestamp:       timestamp,
			MessageLength:   messageLength,
			MessageTypeId:   messageTypeId,
			MessageStreamId: messageStreamId,
		}
	case 1:
		timestampDelta := uint32(uint(message[basicHeader.Length+2]) | uint(message[basicHeader.Length+1])<<8 | uint(message[basicHeader.Length])<<16)
		messageLength := uint32(uint(message[basicHeader.Length+5]) | uint(message[basicHeader.Length+4])<<8 | uint(message[basicHeader.Length+3])<<16)
		messageTypeId := binary.BigEndian.Uint16([]byte{0x00, message[basicHeader.Length+6]})

		return &Header{
			BasicHeader:   basicHeader,
			Timestamp:     timestampDelta,
			MessageLength: messageLength,
			MessageTypeId: messageTypeId,
		}
	case 2:
		timestampDelta := uint32(uint(message[basicHeader.Length+2]) | uint(message[basicHeader.Length+1])<<8 | uint(message[basicHeader.Length])<<16)

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

func getChunkHeaderLength(header *Header) uint32 {
	basicHeaderLength := header.BasicHeader.Length

	switch header.BasicHeader.Type {
	case 0:
		return 11 + basicHeaderLength
	case 1:
		return 7 + 11 + basicHeaderLength
	case 2:
		return 3 + 11 + basicHeaderLength
	default:
		return 0
	}
}
