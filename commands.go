package rtimpus

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type Connect struct {
	CommandName          string
	TransactionID        float64
	CommandObject        Object
	OptionaUserArguments Object
}

type CallResponse struct {
	ProcedureName string
	TransactionID float64
	CommandObject Object
	Response      Object
}

type CreateStream struct {
	CommandName   string
	TransactionID float64
	CommandObject Object
}

const (
	CONNECT       = "connect"
	CREATE_STREAM = "createStream"
)

func UnmarshalCommand(message []byte) (interface{}, error) {
	buffer := bytes.NewBuffer(message)
	decoder := NewAMF0Decoder(buffer)

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

func getConnectResult(decoder *AMF0Decoder) (*Connect, error) {
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

	commObj, ok := commandObject.(Object)

	if !ok {
		return nil, errors.New("commandObject is not of type Object")
	}

	connect.CommandObject = commObj

	optionalUserAgreements, err := decoder.Decode()

	if err != nil && err != io.EOF {
		return nil, err
	}

	optUA, ok := optionalUserAgreements.(Object)

	if !ok {
		connect.OptionaUserArguments = nil
	} else {
		connect.OptionaUserArguments = optUA
	}

	return connect, nil
}

func getCallResponseResult(decoder *AMF0Decoder, value interface{}) (*CallResponse, error) {
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

	commObj, ok := commandObject.(Object)

	if !ok {
		return nil, errors.New("commandObject is not of type Object")
	}

	call.CommandObject = commObj

	response, err := decoder.Decode()

	if err != nil {
		return nil, err
	}

	res, ok := response.(Object)

	if !ok {
		return nil, errors.New("response is not of type Object")
	}

	call.Response = res

	return call, nil
}

func getCreateStream(decoder *AMF0Decoder) (*CallResponse, error) {

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

	commObj, ok := commandObject.(Object)

	if !ok {
		return nil, errors.New("commandObject is not of type Object")
	}

	call.CommandObject = commObj

	return call, nil
}
