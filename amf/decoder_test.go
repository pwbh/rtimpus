package amf

import (
	"bytes"
	"testing"
)

const (
	FIRST_FIELD  = "_result"
	SECOND_FIELD = "Hello World!"
	THIRD_FIELD  = "This is another test"
	FOURTH_FIELD = 10_252
	FIFTH_FIELD  = 50
)

func encodedMessage() ([]byte, error) {
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
	message, err := encodedMessage()
	if err != nil {
		t.Fatalf("coudln't encode a AMF0 message for decoder tests")
	}
	buf := bytes.NewBuffer(message)
	decoder := NewAMF0Decoder(buf)

	firstField, err := decoder.Decode()

	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}

	if firstField.(string) != FIRST_FIELD {
		t.Fatalf("Decoded firstField != FIRST_FIELd")
	}
}
