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

func amf0SimpleMessage() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	amf0Encoder := NewAMF0Encoder(buf)
	if err := amf0Encoder.Encode("_result"); err != nil {
		return nil, err
	}
	if err := amf0Encoder.Encode("Hello World!"); err != nil {
		return nil, err
	}
	if err := amf0Encoder.Encode("This is another test"); err != nil {
		return nil, err
	}
	if err := amf0Encoder.Encode(10_252); err != nil {
		return nil, err
	}
	if err := amf0Encoder.Encode(50); err != nil {
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
