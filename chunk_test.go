package rtimpus

import (
	"fmt"
	"testing"
)

func TestParseChunk(t *testing.T) {
	chunk := parseChunk([]byte{2, 0, 0, 0, 0, 0, 5, 6, 0, 0, 0, 0, 0, 0, 4, 0, 0})

	fmt.Printf("Chunk Type: %d | Chunk Stream ID: %d | Timestamp: %d | Message Length: %d | Message Type ID: %d | Message Stream ID: %d\n", chunk.header.BasicHeader.Type, chunk.header.BasicHeader.StreamID, chunk.header.Timestamp, chunk.header.MessageLength, chunk.header.MessageTypeId, chunk.header.MessageStreamId)
	fmt.Println(chunk.payload.data)
}
