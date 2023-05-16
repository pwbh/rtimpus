package amf

import (
	"bytes"
	"fmt"
	"testing"
)

func TestAMFEncodesMessage(t *testing.T) {
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
	fmt.Println(buf.Bytes())
}
