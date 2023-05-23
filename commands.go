package rtimpus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/pwbh/rtimpus/amf"
	"github.com/pwbh/rtimpus/utils"
)

type Connect struct {
	CommandName          string
	TransactionID        float64
	CommandObject        amf.Object
	OptionaUserArguments amf.Object
}

type CallResponse struct {
	ProcedureName string
	TransactionID float64
	CommandObject amf.Object
	Response      amf.Object
}

type CreateStream struct {
	CommandName   string
	TransactionID float64
	CommandObject amf.Object
}

const (
	CONNECT       = "connect"
	CREATE_STREAM = "createStream"
)

func UnmarshalCommand(chunk *Chunk) (interface{}, error) {
	buffer := bytes.NewBuffer(chunk.payload.data)
	decoder := amf.NewAMF0Decoder(buffer)

	value, err := decoder.Decode()

	if err != nil {
		return nil, err
	}

	switch value {
	case CONNECT:
		return getConnectResult(decoder)
	case CREATE_STREAM:
		return getCreateStream(decoder)

	// If nothing matches then it has to be the incoming RPC call
	default:
		return getCallResponseResult(decoder, value)
	}
}

func getConnectResult(decoder *amf.AMF0Decoder) (*Connect, error) {
	connect := new(Connect)
	connect.CommandName = CONNECT

	transactionID, err := decoder.Decode()

	if err != nil {
		return nil, err
	}

	tranID, ok := transactionID.(float64)

	if !ok {
		return nil, errors.New("transactionID is not of type uint32")
	}

	connect.TransactionID = tranID

	commandObject, err := decoder.Decode()

	if err != nil {
		return nil, err
	}

	commObj, ok := commandObject.(amf.Object)

	if !ok {
		return nil, errors.New("commandObject is not of type Object")
	}

	connect.CommandObject = commObj

	optionalUserAgreements, err := decoder.Decode()

	if err != nil && err != io.EOF {
		return nil, err
	}

	optUA, ok := optionalUserAgreements.(amf.Object)

	if !ok {
		connect.OptionaUserArguments = nil
	} else {
		connect.OptionaUserArguments = optUA
	}

	return connect, nil
}

func getCallResponseResult(decoder *amf.AMF0Decoder, value interface{}) (*CallResponse, error) {
	precedureName, ok := value.(string)

	if !ok {
		return nil, fmt.Errorf("unknown value decoded %v", value)
	}

	call := new(CallResponse)
	call.ProcedureName = precedureName

	transactionID, err := decoder.Decode()

	if err != nil {
		return nil, err
	}

	tranID, ok := transactionID.(float64)

	if !ok {
		return nil, errors.New("transactionID is not of type uint32")
	}

	call.TransactionID = tranID

	commandObject, err := decoder.Decode()

	if err != nil {
		return nil, err
	}

	commObj, ok := commandObject.(amf.Object)

	if !ok {
		return nil, errors.New("commandObject is not of type Object")
	}

	call.CommandObject = commObj

	response, err := decoder.Decode()

	if err != nil {
		return nil, err
	}

	res, ok := response.(amf.Object)

	if !ok {
		return nil, errors.New("response is not of type Object")
	}

	call.Response = res

	return call, nil
}

func getCreateStream(decoder *amf.AMF0Decoder) (*CallResponse, error) {

	call := new(CallResponse)
	call.ProcedureName = CREATE_STREAM

	transactionID, err := decoder.Decode()

	if err != nil {
		return nil, err
	}

	tranID, ok := transactionID.(float64)

	if !ok {
		return nil, errors.New("transactionID is not of type uint32")
	}

	call.TransactionID = tranID

	commandObject, err := decoder.Decode()

	if err != nil {
		return nil, err
	}

	commObj, ok := commandObject.(amf.Object)

	if !ok {
		return nil, errors.New("commandObject is not of type Object")
	}

	call.CommandObject = commObj

	return call, nil
}

// RTMP Chunk Stream uses message type IDs 1, 2, 3, 5, and 6 for protocol control messages.
// These messages contain information needed by the RTMP Chunk Stream protocol.

// These protocol control messages MUST have message stream ID 0 (known as the control stream)
// and be sent in chunk stream ID 2. Protocol control messages take effect as soon as they are received;
// their timestamps are ignored.

// func createProtocolMessageHeader(messageType byte, payloadLength uint32) ([]byte, error) {
// 	if messageType > 6 {
// 		return nil, errors.New("valid messageType ids 1-6, received >6")
// 	}
// 	buf := make([]byte, 12)
// 	buf[0] = 2                              // Chunk Stream ID
// 	buf[1] = messageType                    // Message Type
// 	utils.PutUint24(buf[2:], payloadLength) // Payload length
// 	binary.BigEndian.PutUint32(buf[5:], 0)  // Timestamp is ignored
// 	utils.PutUint24(buf[9:], 0)             // Message Stream ID
// 	return buf, nil
// }

func createProtocolMessageHeader(messageType byte, payloadLength uint32) ([]byte, error) {
	if messageType > 6 {
		return nil, errors.New("valid messageType ids 1-6, received >6")
	}
	buf := make([]byte, 12)
	buf[0] = 2                              // Chunk Stream ID
	utils.PutUint24(buf[1:], 0)             // Timestamp is ignored
	utils.PutUint24(buf[4:], payloadLength) // Payload length
	buf[7] = messageType                    // Message Type
	binary.BigEndian.PutUint32(buf[8:], 0)  // Message Stream ID
	return buf, nil
}

// Protocol control message 1, Set Chunk Size, is used to notify the peer of a new maximum chunk size.
// The maximum chunk size defaults to 128 bytes, but the client or the server can change this value, and updates
// its peer using this message. For example, suppose a client wants to send 131 bytes of audio data and the chunk size is 128.
// In this case, the client can send this message to the server to notify it that the chunk size is now 131 bytes. The client can
// then send the audio data in a single chunk.
// The maximum chunk size SHOULD be at least 128 bytes, and MUST be at least 1 byte. The maximum chunk size
// is maintained independently for each direction.
func sendSetChunkSize(c *Connection, size uint32) error {
	payloadLength := 4
	header, err := createProtocolMessageHeader(1, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, payloadLength+headerLength)
	copy(buf[:headerLength], header)
	binary.BigEndian.PutUint32(buf[headerLength:], size)
	_, wErr := c.Write(buf)
	return wErr
}

// Protocol control message 2, Abort Message, is used to notify the peer if it is waiting for chunks to complete a message,
// then to discard the partially received message over a chunk stream. The peer receives the chunk stream ID as
// this protocol message’s payload. An application may send this message when closing in order to indicate that
// further processing of the messages is not required.
func sendAbortMessage(c *Connection, streamID uint32) error {
	payloadLength := 4
	header, err := createProtocolMessageHeader(2, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, headerLength+payloadLength)
	copy(buf[:headerLength], header)
	binary.BigEndian.PutUint32(buf[headerLength:], streamID)
	_, wErr := c.Write(buf)
	return wErr
}

// Protocol control message 3, Acknowledgement, The client or the server MUST send an acknowledgment to the peer after receiving bytes equal to the window size.
// The window size is the maximum number of bytes that the sender sends without receiving acknowledgment from the receiver.
// This message specifies the sequence number, which is the number of the bytes received so far.
// sequenceNumber field holds the number of bytes received so far.
func sendAcknowledgement(c *Connection, sequenceNumber uint32) error {
	payloadLength := 4
	header, err := createProtocolMessageHeader(3, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, headerLength+payloadLength)
	copy(buf[:headerLength], header)
	binary.BigEndian.PutUint32(buf[headerLength:], sequenceNumber)
	_, wErr := c.Write(buf)
	return wErr
}

// Protocol control message 5, The client or the server sends this message to inform the peer of the window size to use between sending acknowledgments.
// The sender expects acknowledgment from its peer after the sender sends window size bytes.
// The receiving peer MUST send an Acknowledgement (Section 5.4.3) after receiving the indicated
// number of bytes since the last Acknowledgement was sent, or from the beginning of the session if no Acknowledgement has yet been sent.
func sendWindowAcknowledgementSize(c *Connection, size uint32) error {
	payloadLength := 4
	header, err := createProtocolMessageHeader(5, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, headerLength+payloadLength)
	copy(buf[:headerLength], header)
	binary.BigEndian.PutUint32(buf[headerLength:], size)
	_, wErr := c.Write(buf)
	c.WindowSize = size
	return wErr
}

// Protocol contro message 6, The client or the server sends this message to limit the output bandwidth of its peer.
// The peer receiving this message limits its output bandwidth by limiting the amount of sent but unacknowledged
// data to the window size indicated in this message. The peer receiving this message SHOULD
// respond with a Window Acknowledgement Size message if the window size is different from the last
// one sent to the sender of this message.
// The Limit Type is one of the following values:
// 0 - Hard: The peer SHOULD limit its output bandwidth to the indicated window size.
// 1 - Soft: The peer SHOULD limit its output bandwidth to the the window indicated in this message or the limit already in effect, whichever is smaller.
// 2 - Dynamic: If the previous Limit Type was Hard, treat this message as though it was marked Hard, otherwise ignore this message.
func sendSetPeerBandwith(c *Connection, size uint32, limit byte) error {
	if limit > 2 {
		fmt.Printf("limit exceeds maximum of 2, received %d\n", limit)
	}
	payloadLength := 5
	header, err := createProtocolMessageHeader(6, uint32(payloadLength))
	if err != nil {
		return err
	}
	headerLength := len(header)
	buf := make([]byte, headerLength+payloadLength)
	copy(buf[:headerLength], header)
	binary.BigEndian.PutUint32(buf[headerLength:], size)
	buf[headerLength+payloadLength-1] = limit
	_, wErr := c.Write(buf)
	if wErr != nil {
		return wErr
	}
	return nil
}

func sendConnectResult(c *Connection) error {
	buf := bytes.NewBuffer([]byte{})
	encoder := amf.NewAMF0Encoder(buf)
	if err := encoder.Encode("_result"); err != nil {
		return err
	}
	if err := encoder.Encode(1); err != nil {
		return err
	}
	properties := getServerProperties()
	if err := encoder.Encode(properties); err != nil {
		return err
	}
	information := getServerInformation()
	if err := encoder.Encode(information); err != nil {
		return err
	}
	_, err := c.Write(buf.Bytes())
	return err
}

func getServerProperties() amf.Object {
	properties := make(amf.Object)
	properties["fmsVer"] = "FMS/3,5,5,2004"
	properties["capabilities"] = 31.0
	properties["mode"] = 1.0
	return properties
}

func getServerInformation() amf.Object {
	information := make(amf.Object)
	information["level"] = "status"
	information["code"] = "NetConnection.Connect.Success"
	information["description"] = "Connection succeeded"
	return information
}
