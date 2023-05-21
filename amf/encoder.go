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

func (e *AMF0Encoder) encodeArray(ecmaArr ECMAArray) error {
	length := uint32(len(ecmaArr))
	buf := make([]byte, 5)
	buf[0] = AMF0EcmaArrayMarker
	binary.BigEndian.PutUint32(buf[1:], length)
	if _, err := e.writer.Write(buf); err != nil {
		return err
	}
	for _, v := range ecmaArr {
		if err := e.encodeKey(v.Key); err != nil {
			return err
		}
		if err := e.Encode(v.Value); err != nil {
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
	_, err := e.writer.Write(buf)
	return err
}

func (e *AMF0Encoder) encodeKey(str string) error {
	length := len(str)
	buf := make([]byte, 2, length+2)
	binary.BigEndian.PutUint16(buf[0:], uint16(length))
	buf = append(buf, str...)
	_, err := e.writer.Write(buf)
	return err
}

func (e *AMF0Encoder) encodeObject(obj Object) error {
	if _, err := e.writer.Write([]byte{AMF0ObjectMarker}); err != nil {
		return err
	}

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
	if _, err := e.writer.Write([]byte{AMF0NumberMarker, AMF0NumberMarker, AMF0ObjectEndMarker}); err != nil {
		return err
	}

	return nil
}

func (e *AMF0Encoder) encodeNil() error {
	_, err := e.writer.Write([]byte{AMF0NullMarker})
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
