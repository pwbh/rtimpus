package rtimpus

import (
	"fmt"
	"testing"
)

func TestParseChunk(t *testing.T) {
	c := new(Connection)
	chunk, err := parseChunk(c, []byte{2, 0, 0, 0, 0, 0, 5, 6, 0, 0, 0, 0, 0, 0, 4, 0, 0})

	if err != nil {
		t.Fatalf("error while parsing chunk: %v", err)
	}

	fmt.Printf("Chunk Type: %d | Chunk Stream ID: %d | Timestamp: %d | Message Length: %d | Message Type ID: %d | Message Stream ID: %d\n", chunk.header.BasicHeader.Type, chunk.header.BasicHeader.StreamID, chunk.header.MessageHeader.Timestamp, chunk.header.MessageHeader.Length, chunk.header.MessageHeader.TypeID, chunk.header.MessageHeader.StreamID)
	fmt.Println(chunk.payload.data)
}
