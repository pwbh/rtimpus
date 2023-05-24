package rtimpus

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

// The server sends this event to notify the client that a stream has become functional and can be used for communication.
// By default, this event is sent on ID 0 after the application connect command is successfully received from the client.
// The event data is 4-byte and represents the stream ID of the stream that became functional.
func sendStreamBeginEvent(w io.Writer, streamID uint32) error {
	payloadLength := 6
	header, err := createProtocolMessageHeader(4, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, headerLength+payloadLength)
	copy(buf, header)
	binary.BigEndian.PutUint16(buf[headerLength:], 0)
	binary.BigEndian.PutUint32(buf[headerLength+2:], streamID)
	_, wErr := w.Write(buf)
	fmt.Println(buf)
	return wErr
}

// The server sends this event to notify the client that the playback of data is over as requested on this stream.
// No more data is sent without issuing additional commands. The client discards the messages received for the stream.
// The 4 bytes of event data represent the ID of the stream on which playback has ended.
func sendStreamEOFEvent(w io.Writer, streamID uint32) error {
	payloadLength := 6
	header, err := createProtocolMessageHeader(4, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, headerLength+payloadLength)
	copy(buf, header)
	binary.BigEndian.PutUint16(buf[headerLength:], 1)
	binary.BigEndian.PutUint32(buf[headerLength+2:], streamID)
	_, wErr := w.Write(buf)
	return wErr
}

// The server sends this event to notify the client that there is no more data on the stream.
// If the server does not detect any message for a time period, it can notify the subscribed clients that the stream is dry.
// The 4 bytes of event data represent the stream ID of the dry stream.
func sendStreamDryEvent(w io.Writer, streamID uint32) error {
	payloadLength := 6
	header, err := createProtocolMessageHeader(4, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, headerLength+payloadLength)
	copy(buf, header)
	binary.BigEndian.PutUint16(buf[headerLength:], 2)
	binary.BigEndian.PutUint32(buf[headerLength+2:], streamID)
	_, wErr := w.Write(buf)
	return wErr
}

// The client sends this event to inform the server of the buffer size (in milliseconds) that is used to buffer any data coming over a stream.
// This event is sent before the server starts processing the stream.
// The first 4 bytes of the event data represent the stream ID and the next 4 bytes represent the buffer length, in milliseconds.
func sendBufferLengthEvent(w io.Writer, streamID uint32, bufferLength uint32) error {
	payloadLength := 10
	header, err := createProtocolMessageHeader(4, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, headerLength+payloadLength)
	copy(buf, header)
	binary.BigEndian.PutUint16(buf[headerLength:], 3)
	binary.BigEndian.PutUint32(buf[headerLength+2:], streamID)
	binary.BigEndian.PutUint32(buf[headerLength+6:], bufferLength)
	_, wErr := w.Write(buf)
	return wErr
}

// The server sends this event to notify the client that the stream is a recorded stream.
// The 4 bytes event data represent the stream ID of the recorded stream.
func sendStreamlsRecordedEvent(w io.Writer, streamID uint32) error {
	payloadLength := 6
	header, err := createProtocolMessageHeader(4, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, headerLength+payloadLength)
	copy(buf, header)
	binary.BigEndian.PutUint16(buf[headerLength:], 4)
	binary.BigEndian.PutUint32(buf[headerLength+2:], streamID)
	_, wErr := w.Write(buf)
	return wErr
}

// The server sends this event to test whether the client is reachable.
// Event data is a 4-byte timestamp, representing the local server time when the server dispatched the command.
// The client responds with PingResponse on receiving MsgPingRequest.
func sendPingRequestEvent(w io.Writer) error {
	payloadLength := 6
	header, err := createProtocolMessageHeader(4, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, headerLength+payloadLength)
	copy(buf, header)
	binary.BigEndian.PutUint16(buf[headerLength:], 6)
	binary.BigEndian.PutUint32(buf[headerLength+2:], uint32(time.Now().Unix()))
	_, wErr := w.Write(buf)
	return wErr
}

// The client sends this event to the server in response to the ping request.
// The event data is a 4-byte timestamp, which was received with the PingRequest request.
func sendPingResponseEvent(w io.Writer, timestamp uint32) error {
	payloadLength := 6
	header, err := createProtocolMessageHeader(4, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, headerLength+payloadLength)
	copy(buf, header)
	binary.BigEndian.PutUint16(buf[headerLength:], 7)
	binary.BigEndian.PutUint32(buf[headerLength+2:], timestamp)
	_, wErr := w.Write(buf)
	return wErr
}
