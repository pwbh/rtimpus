package amf

import (
	"bytes"
	"testing"
)

const (
	SIMPLE_MESSAGE_FIRST_FIELD  = "_result"
	SIMPLE_MESSAGE_SECOND_FIELD = "Hello World!"
	SIMPLE_MESSAGE_THIRD_FIELD  = "This is another test"
	SIMPLE_MESSAGE_FOURTH_FIELD = 10_252
	SIMPLE_MESSAGE_FIFTH_FIELD  = 50.0003
)

const (
	COMPLEX_MESSAGE_FIRST_FIELD  = "_result"
	COMPLEX_MESSAGE_SECOND_FIELD = "version"
	COMPLEX_MESSAGE_THIRD_FIELD  = 32653.12414
	COMPLEX_MESSAGE_FOURTH_FIELD = "other_ObJeCt"

	COMPLEX_MESSAGE_OBJECT_FIRST_KEY    = "scvWS"
	COMPLEX_MESSAGE_OBJECT_FIRST_VALUE  = "https://localhost:10000/"
	COMPLEX_MESSAGE_OBJECT_SECOND_KEY   = "test"
	COMPLEX_MESSAGE_OBJECT_SECOND_VALUE = "test"
	COMPLEX_MESSAGE_OBJECT_THIRD_KEY    = "another_test"
	COMPLEX_MESSAGE_OBJECT_THIRD_VALUE  = 50.1231

	COMPLEX_MESSAGE_OBJECT_FOURTH_KEY = "innerObj"

	COMPLEX_MESSAGE_OBJECT_FOURTH_KEY_INNER_OBJECT_FIRST_KEY   = "test"
	COMPLEX_MESSAGE_OBJECT_FOURTH_KEY_INNER_OBJECT_FIRST_VALUE = "test"

	COMPLEX_MESSAGE_FIFTH_KEY   = "otherKey"
	COMPLEX_MESSAGE_FIFTH_VALUE = 2
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
	if err := amf0Encoder.Encode(COMPLEX_MESSAGE_THIRD_FIELD); err != nil {
		return nil, err
	}
	complexObject := make(Object)
	complexObject[COMPLEX_MESSAGE_OBJECT_FIRST_KEY] = COMPLEX_MESSAGE_OBJECT_FIRST_VALUE
	complexObject[COMPLEX_MESSAGE_OBJECT_SECOND_KEY] = COMPLEX_MESSAGE_OBJECT_SECOND_VALUE
	complexObject[COMPLEX_MESSAGE_OBJECT_THIRD_KEY] = COMPLEX_MESSAGE_OBJECT_THIRD_VALUE
	innerObject := make(Object)
	innerObject[COMPLEX_MESSAGE_OBJECT_FOURTH_KEY_INNER_OBJECT_FIRST_KEY] = COMPLEX_MESSAGE_OBJECT_FOURTH_KEY_INNER_OBJECT_FIRST_VALUE
	complexObject[COMPLEX_MESSAGE_OBJECT_FOURTH_KEY] = innerObject
	complexObject[COMPLEX_MESSAGE_FIFTH_KEY] = COMPLEX_MESSAGE_FIFTH_VALUE
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
	if fourthField.(float64) != float64(SIMPLE_MESSAGE_FOURTH_FIELD) {
		t.Fatalf("Decoded fourthField != SIMPLE_MESSAGE_FOURTH_FIELD")
	}
	fifthField, err := decoder.Decode()
	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}
	if fifthField.(float64) != float64(SIMPLE_MESSAGE_FIFTH_FIELD) {
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
	secondField, err := decoder.Decode()
	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}
	if secondField.(string) != COMPLEX_MESSAGE_SECOND_FIELD {
		t.Fatalf("Decoded secondField != COMPLEX_MESSAGE_SECOND_FIELD")
	}
	thirdField, err := decoder.Decode()
	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}
	if thirdField.(float64) != COMPLEX_MESSAGE_THIRD_FIELD {
		t.Fatalf("Decoded thirdField != COMPLEX_MESSAGE_THIRD_FIELD")
	}
	fourthField, err := decoder.Decode()
	if err != nil {
		t.Fatalf("err while decoding: %v", err)
	}
	if v, ok := fourthField.(Object); ok {
		if v[COMPLEX_MESSAGE_OBJECT_FIRST_KEY].(string) != COMPLEX_MESSAGE_OBJECT_FIRST_VALUE {
			t.Fatalf(`Decoded v[COMPLEX_MESSAGE_THIRD_FIELD_OBJECT_FIRST_KEY] != COMPLEX_MESSAGE_OBJECT_FIRST_VALUE`)
		}
		if v[COMPLEX_MESSAGE_OBJECT_SECOND_KEY].(string) != COMPLEX_MESSAGE_OBJECT_SECOND_VALUE {
			t.Fatalf(`Decoded v[COMPLEX_MESSAGE_OBJECT_SECOND_KEY] != COMPLEX_MESSAGE_OBJECT_SECOND_VALUE`)
		}
		if v[COMPLEX_MESSAGE_OBJECT_THIRD_KEY].(float64) != COMPLEX_MESSAGE_OBJECT_THIRD_VALUE {
			t.Fatalf(`Decoded v[COMPLEX_MESSAGE_OBJECT_THIRD_KEY] != COMPLEX_MESSAGE_OBJECT_THIRD_VALUE`)
		}
		if innerObj, ok := v[COMPLEX_MESSAGE_OBJECT_FOURTH_KEY].(Object); ok {
			if innerObj[COMPLEX_MESSAGE_OBJECT_FOURTH_KEY_INNER_OBJECT_FIRST_KEY].(string) != COMPLEX_MESSAGE_OBJECT_FOURTH_KEY_INNER_OBJECT_FIRST_VALUE {
				t.Fatalf(`Decoded innerObj[COMPLEX_MESSAGE_OBJECT_FOURTH_KEY_INNER_OBJECT_FIRST_KEY] != COMPLEX_MESSAGE_OBJECT_FOURTH_KEY_INNER_OBJECT_FIRST_VALUE`)
			}
		} else {
			t.Fatalf(`Type of v[COMPLEX_MESSAGE_OBJECT_FOURTH_KEY] != Object`)
		}
		if v[COMPLEX_MESSAGE_FIFTH_KEY].(float64) != COMPLEX_MESSAGE_FIFTH_VALUE {
			t.Fatalf(`Decoded v[COMPLEX_MESSAGE_FIFTH_KEY] != COMPLEX_MESSAGE_FIFTH_VALUE`)
		}
	} else {
		t.Fatalf("Decoded thirdField != Object")
	}
}
