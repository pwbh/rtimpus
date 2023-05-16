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
	buf := make([]byte, 9)
	buf[0] = AMF0NumberMarker
	binary.BigEndian.PutUint64(buf[1:], math.Float64bits(num))
	_, err := e.writer.Write(buf)
	return err
}

func (e *AMF0Encoder) encodeString(str string) error {
	length := len(str)
	buf := make([]byte, 3, length+3)
	buf[0] = AMF0StringMarker
	binary.BigEndian.PutUint16(buf[1:], uint16(length))
	buf = append(buf, str...)
	_, err := e.writer.Write(buf)
	return err
}

func (e *AMF0Encoder) encodeArray(arr []uint32) error {
	length := len(arr) * 4
	buf := make([]byte, length+1)
	buf[0] = AMF0EcmaArrayMarker
	for i, e := range arr {
		binary.BigEndian.AppendUint32(buf[i*4+1:i*4+5], e)
	}
	return nil
}

func (e *AMF0Encoder) encodeBool(v bool) error {
	value := byte(0)
	if v {
		value = 1
	}
	buf := []byte{AMF0BooleanMarker, value}
	_, err := e.writer.Write(buf)
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
		case bool:
			if err := e.encodeBool(v); err != nil {
				return err
			}
		case Object:
			if err := e.encodeObject(v); err != nil {
				return err
			}
		case string:
			if err := e.encodeString(v); err != nil {
				return err
			}
		case int:
			if err := e.encodeNumber(float64(v)); err != nil {
				return err
			}
		case float32:
			if err := e.encodeNumber(float64(v)); err != nil {
				return err
			}
		case float64:
			if err := e.encodeNumber(v); err != nil {
				return err
			}
		case []uint32: // Probably will not work correctly yet
			if err := e.encodeArray(v); err != nil {
				return err
			}
		default:
			return errors.New("type is not recognized")
		}
	}

	// As noted in the specification - 0x00 0x00 0x09
	if _, err := e.writer.Write([]byte{AMF0NumberMarker, AMF0NumberMarker, AMF0ObjectEndMarker}); err != nil {
		return err
	}

	return nil
}

func (e *AMF0Encoder) Encode(value interface{}) error {
	switch v := value.(type) {
	case bool:
		if err := e.encodeBool(v); err != nil {
			return err
		}
	case Object:
		if err := e.encodeObject(v); err != nil {
			return err
		}
	case string:
		if err := e.encodeString(v); err != nil {
			return err
		}
	case int:
		if err := e.encodeNumber(float64(v)); err != nil {
			return err
		}
	case float32:
		if err := e.encodeNumber(float64(v)); err != nil {
			return err
		}
	case float64:
		if err := e.encodeNumber(v); err != nil {
			return err
		}
	case []uint32: // Probably will not work correctly yet
		if err := e.encodeArray(v); err != nil {
			return err
		}
	default:
		return errors.New("type is not recognized")
	}
	return nil
}
