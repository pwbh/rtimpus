package amf

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

type AMF0Encoder struct {
	writer io.Writer
}

func NewAMF0Encoder(writer io.Writer) *AMF0Encoder {
	return &AMF0Encoder{writer: writer}
}

func (e *AMF0Encoder) encodeNumber(num float64) error {
	data := make([]byte, 1, 9)
	data[0] = AMF0NumberMarker
	binary.BigEndian.PutUint64(data[1:], math.Float64bits(num))
	_, err := e.writer.Write(data)
	return err
}

func (e *AMF0Encoder) encodeString(str string) error {
	length := len(str)
	data := make([]byte, 3, length+3)
	data[0] = AMF0StringMarker
	binary.BigEndian.PutUint16(data[:1], uint16(length))
	data = append(data, str...)
	_, err := e.writer.Write(data)
	return err
}

func (e *AMF0Encoder) encodeBool(v bool) error {
	value := byte(0)
	if v {
		value = 1
	}
	data := []byte{AMF0BooleanMarker, value}
	_, err := e.writer.Write(data)
	return err
}

func (e *AMF0Encoder) encodeObject(obj Object) error {
	if _, err := e.writer.Write([]byte{AMF0ObjectMarker}); err != nil {
		return err
	}

	for k, v := range obj {
		if err := e.encodeString(k); err != nil {
			return err
		}
		switch v := v.(type) {
		case Object:
			if err := e.encodeObject(v); err != nil {
				return err
			}
		case string:
			if err := e.encodeString(v); err != nil {
				return err
			}
		case float64:
			if err := e.encodeNumber(v); err != nil {
				return err
			}
		case []uint32:
			// need to encode arrays once encodeArray is ready
		default:
			return errors.New("type is not recognized")
		}
	}

	if _, err := e.writer.Write([]byte{AMF0ObjectEndMarker}); err != nil {
		return err
	}

	return nil
}
