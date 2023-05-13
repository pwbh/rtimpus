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

func UnmarshalCommand(message []byte) (interface{}, error) {
	buffer := bytes.NewBuffer(message)
	decoder := NewAMF0Decoder(buffer)

	command, err := decoder.Decode()

	fmt.Println(command)

	if err != nil {
		return nil, err
	}

	switch command {
	case "connect":
		return getConnectResult(decoder)

	default:
		return nil, fmt.Errorf("no such command has been found, %s", command)
	}
}

func getConnectResult(decoder *AMF0Decoder) (*Connect, error) {
	connect := new(Connect)
	connect.CommandName = "connect"

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
