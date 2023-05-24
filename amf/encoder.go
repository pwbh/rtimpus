package amf

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

type AMF0Encoder struct {
	writer io.Writer
	length int
}

func NewAMF0Encoder(writer io.Writer) *AMF0Encoder {
	return &AMF0Encoder{writer: writer}
}

func (e *AMF0Encoder) encodeNumber(num float64) error {
	buf := make([]byte, 9)
	buf[0] = AMF0NumberMarker
	binary.BigEndian.PutUint64(buf[1:], math.Float64bits(num))
	n, err := e.writer.Write(buf)
	e.length += n
	return err
}

func (e *AMF0Encoder) encodeString(str string) error {
	length := len(str)
	buf := make([]byte, 3, length+3)
	buf[0] = AMF0StringMarker
	binary.BigEndian.PutUint16(buf[1:], uint16(length))
	buf = append(buf, str...)
	n, err := e.writer.Write(buf)
	e.length += n
	return err
}

func (e *AMF0Encoder) encodeArray(ecmaArr ECMAArray) error {
	length := uint32(len(ecmaArr))
	buf := make([]byte, 5)
	buf[0] = AMF0EcmaArrayMarker
	binary.BigEndian.PutUint32(buf[1:], length)
	n, err := e.writer.Write(buf)
	if err != nil {
		return err
	}
	e.length += n
	for _, v := range ecmaArr {

		if err := e.encodeKey(v.K); err != nil {
			return err
		}
		if err := e.Encode(v.V); err != nil {
			return err
		}
	}
	return nil
}

func (e *AMF0Encoder) encodeBool(v bool) error {
	value := byte(0)
	if v {
		value = 1
	}
	buf := []byte{AMF0BooleanMarker, value}
	n, err := e.writer.Write(buf)
	e.length += n
	return err
}

func (e *AMF0Encoder) encodeKey(str string) error {
	length := len(str)
	buf := make([]byte, 2, length+2)
	binary.BigEndian.PutUint16(buf[0:], uint16(length))
	buf = append(buf, str...)
	n, err := e.writer.Write(buf)
	e.length += n
	return err
}

func (e *AMF0Encoder) encodeObject(obj Object) error {
	n, err := e.writer.Write([]byte{AMF0ObjectMarker})
	if err != nil {
		return err
	}
	e.length += n
	for k, v := range obj {
		if err := e.encodeKey(k); err != nil {
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
		case ECMAArray:
			if err := e.encodeArray(v); err != nil {
				return err
			}
		default:
			return errors.New("type is not recognized")
		}
	}
	// As noted in the specification - 0x00 0x00 0x09
	k, err := e.writer.Write([]byte{AMF0NumberMarker, AMF0NumberMarker, AMF0ObjectEndMarker})
	if err != nil {
		return err
	}
	e.length += k
	return nil
}

func (e *AMF0Encoder) encodeNil() error {
	n, err := e.writer.Write([]byte{AMF0NullMarker})
	e.length += n
	return err
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
	case ECMAArray:
		if err := e.encodeArray(v); err != nil {
			return err
		}
	case nil:
		if err := e.encodeNil(); err != nil {
			return err
		}
	default:
		return errors.New("type is not recognized")
	}
	return nil
}

func (e *AMF0Encoder) Length() int {
	return e.length
}
