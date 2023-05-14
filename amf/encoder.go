package amf

import (
	"encoding/binary"
	"io"
)

type AMF0Encoder struct {
	writer io.Writer
}

func NewAMF0Encoder(writer io.Writer) *AMF0Encoder {
	return &AMF0Encoder{writer: writer}
}

func (e *AMF0Encoder) encodeNumber(num float64) error {
	return binary.Write(e.writer, binary.BigEndian, num)
}

func (e *AMF0Encoder) encodeString(str string) error {
	length := len(str)

	data := make([]byte, length+1)
	data[0] = byte(length)
	data = append(data, str...)

	_, err := e.writer.Write(data)

	return err
}
