package amf

import (
	"bytes"
	"testing"
)

func TestAMF0EncoderBasicMessage(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	amf0Encoder := NewAMF0Encoder(buf)
	if err := amf0Encoder.Encode("_result"); err != nil {
		t.Fatalf(`amf0Encoder.Encode("_result") =  %v`, err)
	}
	if err := amf0Encoder.Encode("Hello World!"); err != nil {
		t.Fatalf(`amf0Encoder.Encode("Hello World!") =  %v`, err)
	}
	if err := amf0Encoder.Encode("This is another test"); err != nil {
		t.Fatalf(`amf0Encoder.Encode("This is another test") =  %v`, err)
	}
	if err := amf0Encoder.Encode(10_252); err != nil {
		t.Fatalf(`amf0Encoder.Encode("This is another test") =  %v`, err)
	}
	if err := amf0Encoder.Encode(50); err != nil {
		t.Fatalf(`amf0Encoder.Encode("This is another test") =  %v`, err)
	}
}

func TestAMF0EncoderComplexMessage(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	amf0Encoder := NewAMF0Encoder(buf)
	if err := amf0Encoder.Encode("_result"); err != nil {
		t.Fatalf(`amf0Encoder.Encode("_result") =  %v`, err)
	}
	if err := amf0Encoder.Encode(1); err != nil {
		t.Fatalf(`amf0Encoder.Encode(1) =  %v`, err)
	}
	complexObject := Object{"version": 3, "something": "whatever string we want", "swVcs": "https://localhost:10000/", "some_else": Object{"test": "test"}}
	if err := amf0Encoder.Encode(complexObject); err != nil {
		t.Fatalf(`amf0Encoder.Encode(complexObject) =  %v`, err)
	}
}
