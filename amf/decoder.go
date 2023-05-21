package amf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type AMF0Decoder struct {
	reader io.Reader
}

func NewAMF0Decoder(reader io.Reader) *AMF0Decoder {
	return &AMF0Decoder{reader: reader}
}

func (d *AMF0Decoder) readByte() (byte, error) {
	var b [1]byte
	_, err := d.reader.Read(b[:])
	return b[0], err
}

func (d *AMF0Decoder) readBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := d.reader.Read(b)
	return b, err
}

func (d *AMF0Decoder) readUint16() (uint16, error) {
	var b [2]byte
	_, err := d.reader.Read(b[:])
	return binary.BigEndian.Uint16(b[:]), err
}

func (d *AMF0Decoder) readUint32() (uint32, error) {
	var b [4]byte
	_, err := d.reader.Read(b[:])
	return binary.BigEndian.Uint32(b[:]), err
}

func (d *AMF0Decoder) Decode() (interface{}, error) {
	marker, err := d.readByte()
	if err != nil {
		return nil, err
	}
	switch marker {
	case AMF0NumberMarker:
		return d.decodeNumber()
	case AMF0BooleanMarker:
		return d.decodeBoolean()
	case AMF0StringMarker:
		return d.decodeString()
	case AMF0ObjectMarker:
		return d.decodeObject()
	case AMF0NullMarker:
		return nil, nil
	case AMF0UndefinedMarker:
		return nil, nil
	case AMF0ReferenceMarker:
		return nil, errors.New("AMF0 reference not supported")
	case AMF0EcmaArrayMarker:
		return d.decodeEcmaArray()
	case AMF0ObjectEndMarker:
		return nil, errors.New("unexpected AMF0 object end marker")
	default:
		return nil, fmt.Errorf("unknown AMF0 marker: %02x", marker)
	}
}

func (d *AMF0Decoder) decodeNumber() (float64, error) {
	var num float64
	err := binary.Read(d.reader, binary.BigEndian, &num)
	return num, err
}

func (d *AMF0Decoder) decodeBoolean() (bool, error) {
	b, err := d.readByte()
	if err != nil {
		return false, err
	}
	return b != 0, nil
}

func (d *AMF0Decoder) decodeString() (string, error) {
	length, err := d.readUint16()
	if err != nil {
		return "", err
	}
	data, err := d.readBytes(int(length))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (d *AMF0Decoder) decodeObject() (Object, error) {
	obj := make(Object)
	for {
		name, err := d.decodeString()
		if err != nil {
			return nil, err
		}
		if name == "" {
			marker, err := d.readByte()

			if err != nil {
				return nil, err
			}

			if marker != AMF0ObjectEndMarker {
				return nil, errors.New("missing AMF0 object end marker")
			}

			break
		}
		value, err := d.Decode()
		if err != nil {
			return nil, err
		}
		obj[name] = value
	}
	return obj, nil
}

func (d *AMF0Decoder) decodeEcmaArray() (ECMAArray, error) {
	length, err := d.readUint32()
	if err != nil {
		return nil, err
	}
	arr := make(ECMAArray, length)
	for i := uint32(0); i < length; i++ {
		_, err := d.decodeString()
		if err != nil {
			return nil, err
		}
		value, err := d.Decode()
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)
	}
	return arr, nil
}
