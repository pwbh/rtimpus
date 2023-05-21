package rtimpus

import (
	"encoding/binary"
	"io"
)

// The server sends this event to notify the client that a stream has become functional and can be used for communication.
// By default, this event is sent on ID 0 after the application connect command is successfully received from the client.
// The event data is 4-byte and represents the stream ID of the stream that became functional.
func SendStreamBeginEvent(w io.Writer, streamID uint32) error {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[:2], 0)
	binary.BigEndian.PutUint32(buf[2:], streamID)
	_, err := w.Write(buf)
	return err
}
