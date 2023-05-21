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
	complexObject := Object{"version": 3, "something": "whatever string we want", "swVcs": "https://localhost:10000/", "some_else": Object{"test": "test"}, "sometrhing": 2.0}
	if err := amf0Encoder.Encode(complexObject); err != nil {
		t.Fatalf(`amf0Encoder.Encode(complexObject) =  %v`, err)
	}
}

func TestAMF0EncoderEncodeECMAArray(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	amf0Encoder := NewAMF0Encoder(buf)
	if err := amf0Encoder.Encode("_result"); err != nil {
		t.Fatalf(`amf0Encoder.Encode("_result") =  %v`, err)
	}
	if err := amf0Encoder.Encode(1); err != nil {
		t.Fatalf(`amf0Encoder.Encode(1) =  %v`, err)
	}
	complexObject := Object{"version": 3, "something": "whatever string we want", "swVcs": "https://localhost:10000/", "some_else": Object{"test": "test"}, "sometrhing": 2.0}
	if err := amf0Encoder.Encode(complexObject); err != nil {
		t.Fatalf(`amf0Encoder.Encode(complexObject) =  %v`, err)
	}
	ecmaArr := make(ECMAArray, 0, 5)
	ecmaArr = append(ecmaArr, ECMAArrayItem{K: "Name", V: "World"})
	ecmaArr = append(ecmaArr, ECMAArrayItem{K: "version", V: 2.1})
	ecmaArr = append(ecmaArr, ECMAArrayItem{K: "url", V: "https://localhost:10000/"})
	ecmaArr = append(ecmaArr, ECMAArrayItem{K: "inner", V: ECMAArray{ECMAArrayItem{K: "something", V: "Test"}}})
	ecmaArr = append(ecmaArr, ECMAArrayItem{K: "lastItem", V: "World2"})
	if err := amf0Encoder.Encode(ecmaArr); err != nil {
		t.Fatalf(`amf0Encoder.Encode(ecmaArr) =  %v`, err)
	}
}

func BenchmarkAMF0EncoderComplexType(b *testing.B) {
	buf := bytes.NewBuffer([]byte{})
	amf0Encoder := NewAMF0Encoder(buf)

	complexObject := Object{
		"version":    3,
		"something":  "whatever string we want",
		"swVcs":      "https://localhost:10000/",
		"some_else":  Object{"test": "test"},
		"sometrhing": 2.0}

	if err := amf0Encoder.Encode("_result"); err != nil {
		b.Fatalf(`amf0Encoder.Encode("_result") =  %v`, err)
	}
	if err := amf0Encoder.Encode("Hello World!"); err != nil {
		b.Fatalf(`amf0Encoder.Encode("Hello World!") =  %v`, err)
	}
	if err := amf0Encoder.Encode(complexObject); err != nil {
		b.Fatalf(`amf0Encoder.Encode(complexObject) =  %v`, err)
	}
	if err := amf0Encoder.Encode("This is another test"); err != nil {
		b.Fatalf(`amf0Encoder.Encode("This is another test") =  %v`, err)
	}
	if err := amf0Encoder.Encode(10_252); err != nil {
		b.Fatalf(`amf0Encoder.Encode("This is another test") =  %v`, err)
	}
	if err := amf0Encoder.Encode(50); err != nil {
		b.Fatalf(`amf0Encoder.Encode("This is another test") =  %v`, err)
	}
}
