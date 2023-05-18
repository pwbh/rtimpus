package amf

import (
	"bytes"
	"testing"
)

const (
	SIMPLE_MESSAGE_FIRST_FIELD  = "_result"
	SIMPLE_MESSAGE_SECOND_FIELD = "Hello World!"
	SIMPLE_MESSAGE_THIRD_FIELD  = "This is another test"
	SIMPLE_MESSAGE_FOURTH_FIELD = float64(10_252)
	SIMPLE_MESSAGE_FIFTH_FIELD  = float64(50)
)

const (
	COMPLEX_MESSAGE_FIRST_FIELD  = "_result"
	COMPLEX_MESSAGE_SECOND_FIELD = float64(1)

	COMPLEX_MESSAGE_THIRD_FIELD_OBJECT_FIRST_VALUE                     = 3
	COMPLEX_MESSAGE_THIRD_FIELD_OBJECT_SECOND_VALUE                    = "whatever string we want"
	COMPLEX_MESSAGE_THIRD_FIELD_OBJECT_THIRD_VALUE                     = "https://localhost:10000/"
	COMPLEX_MESSAGE_THIRD_FIELD_OBJECT_FOURTH_VALUE_OBJECT_FIRST_VALUE = "test"
	COMPLEX_MESSAGE_THIRD_FIELD_OBJECT_FIFTH_VALUE                     = float64(2.0)
)

func amf0SimpleMessage() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	amf0Encoder := NewAMF0Encoder(buf)
	if err := amf0Encoder.Encode(SIMPLE_MESSAGE_FIRST_FIELD); err != nil {
		return nil, err
	}
	if err := amf0Encoder.Encode(SIMPLE_MESSAGE_SECOND_FIELD); err != nil {
		return nil, err
	}
	if err := amf0Encoder.Encode(SIMPLE_MESSAGE_THIRD_FIELD); err != nil {
		return nil, err
	}
	if err := amf0Encoder.Encode(SIMPLE_MESSAGE_FOURTH_FIELD); err != nil {
		return nil, err
	}
	if err := amf0Encoder.Encode(SIMPLE_MESSAGE_FIFTH_FIELD); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func amf0ComplexMessage() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	amf0Encoder := NewAMF0Encoder(buf)
	if err := amf0Encoder.Encode(COMPLEX_MESSAGE_FIRST_FIELD); err != nil {
		return nil, err
	}
	if err := amf0Encoder.Encode(COMPLEX_MESSAGE_SECOND_FIELD); err != nil {
		return nil, err
	}
	complexObject := Object{
		"version":    COMPLEX_MESSAGE_THIRD_FIELD_OBJECT_FIRST_VALUE,
		"something":  COMPLEX_MESSAGE_THIRD_FIELD_OBJECT_SECOND_VALUE,
		"swVcs":      COMPLEX_MESSAGE_THIRD_FIELD_OBJECT_THIRD_VALUE,
		"some_else":  Object{"test": COMPLEX_MESSAGE_THIRD_FIELD_OBJECT_FOURTH_VALUE_OBJECT_FIRST_VALUE},
		"sometrhing": COMPLEX_MESSAGE_THIRD_FIELD_OBJECT_FIFTH_VALUE}
	if err := amf0Encoder.Encode(complexObject); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func TestAMF0DecoderDecodesSimpleMessage(t *testing.T) {
	message, err := amf0SimpleMessage()
	if err != nil {
		t.Fatalf("coudln't encode a AMF0 message for decoder tests")
	}
	buf := bytes.NewBuffer(message)
	decoder := NewAMF0Decoder(buf)
	firstField, err := decoder.Decode()
	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}
	if firstField.(string) != SIMPLE_MESSAGE_FIRST_FIELD {
		t.Fatalf("Decoded firstField != SIMPLE_MESSAGE_FIRST_FIELD")
	}
	secondField, err := decoder.Decode()
	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}
	if secondField.(string) != SIMPLE_MESSAGE_SECOND_FIELD {
		t.Fatalf("Decoded secondField != SIMPLE_MESSAGE_SECOND_FIELD")
	}
	thirdField, err := decoder.Decode()
	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}
	if thirdField.(string) != SIMPLE_MESSAGE_THIRD_FIELD {
		t.Fatalf("Decoded thirdField != SIMPLE_MESSAGE_THIRD_FIELD")
	}
	fourthField, err := decoder.Decode()
	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}
	if fourthField.(float64) != SIMPLE_MESSAGE_FOURTH_FIELD {
		t.Fatalf("Decoded fourthField != SIMPLE_MESSAGE_FOURTH_FIELD")
	}
	fifthField, err := decoder.Decode()
	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}
	if fifthField.(float64) != SIMPLE_MESSAGE_FIFTH_FIELD {
		t.Fatalf("Decoded fifthField != SIMPLE_MESSAGE_FIFTH_FIELD")
	}
}

func TestAMF0DecoderDecodesComplexMessage(t *testing.T) {
	message, err := amf0ComplexMessage()
	if err != nil {
		t.Fatalf("coudln't encode a AMF0 message for decoder tests")
	}
	buf := bytes.NewBuffer(message)
	decoder := NewAMF0Decoder(buf)
	firstField, err := decoder.Decode()
	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}
	if firstField.(string) != COMPLEX_MESSAGE_FIRST_FIELD {
		t.Fatalf("Decoded firstField != COMPLEX_MESSAGE_FIRST_FIELD")
	}
	s, err := decoder.Decode()
	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}
	if s.(float64) != COMPLEX_MESSAGE_SECOND_FIELD {
		t.Fatalf("Decoded s != COMPLEX_MESSAGE_SECOND_FIELD")
	}
}
