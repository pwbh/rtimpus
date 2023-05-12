package rtimpus

import (
	"encoding/binary"
	"fmt"
)

type AMF0 byte

func (a AMF0) String() string {
	switch a {
	case AMF0Number:
		return "AMF0Number"
	case AMF0Boolean:
		return "AMF0Boolean"
	case AMF0String:
		return "AMF0String"
	case AMF0Object:
		return "AMF0Object"
	case AMF0Null:
		return "AMF0Null"
	case AMF0ECMAArray:
		return "AMF0ECMAArray"
	case AMF0ObjectEnd:
		return "AMF0ObjectEnd"
	case AMF0StrictArray:
		return "AMF0StrictArray"
	case AMF0Date:
		return "AMF0Date"
	case AMF0LongString:
		return "AMF0LongString"
	case AMF0XMLDocument:
		return "AMF0XMLDocument"
	case AMF0TypedObject:
		return "AMF0TypedObject"
	case AMF0SwitchToAMF3:
		return "AMF0SwitchToAMF3"
	default:
		return "error: unrecognized AMF0 type"
	}
}

const (
	AMF0Number       = AMF0(0x00) // Number
	AMF0Boolean      = AMF0(0x01) // Boolean
	AMF0String       = AMF0(0x02) // String
	AMF0Object       = AMF0(0x03) // Object (Set of key/value pairs)
	AMF0Null         = AMF0(0x05) // Null
	AMF0ECMAArray    = AMF0(0x08) // ECMA Array
	AMF0ObjectEnd    = AMF0(0x09) // Object End
	AMF0StrictArray  = AMF0(0x0a) // Strict Array
	AMF0Date         = AMF0(0x0b) // Date
	AMF0LongString   = AMF0(0x0c) // Long String
	AMF0XMLDocument  = AMF0(0x0f) // XML Document
	AMF0TypedObject  = AMF0(0x10) // Typed Object
	AMF0SwitchToAMF3 = AMF0(0x11) // Switch to AMF3
)

type Value struct {
	Type AMF0
	Data []byte
}

type AMF0Result struct {
	Command       string
	TransactionID uint16
	Payload       []Value
}

func UnmarshalAMF0(message []byte) *AMF0Result {
	result := new(AMF0Result)

	phase := 0
	start := 0
	end := 0

	for len(message) > end {
		offset, amf0Type, ok := getFoundOffset(message, end)

		if ok && end > start {
			if phase == 0 {
				result.Command = string(message[start:end])
				phase = 1
			} else if phase == 1 {
				result.TransactionID = binary.BigEndian.Uint16(message[start:end])
				phase = 2
			} else {
				fmt.Println(offset, amf0Type, string(message[start:end]), message[start:end])
			}

			start = end
		}

		if offset > 0 {
			start += offset
			end += offset + 1
		} else {
			end++
		}
	}

	return result
}

func getFoundOffset(message []byte, currentIndex int) (int, AMF0, bool) {
	a, ok := getAMFFieldType(message[currentIndex])

	if !ok {
		return 0, a, false
	}

	b := false

	if len(message) > currentIndex+1 {
		b = message[currentIndex+1] == byte(AMF0Number)
	}

	if ok && b {
		if a == AMF0Number {
			return 0, a, false
		}

		return 2, a, true
	}

	return 1, a, true
}

func getAMFFieldType(b byte) (AMF0, bool) {
	AMF0Types := [...]AMF0{
		AMF0Number,
		AMF0Boolean,
		AMF0String,
		AMF0Object,
		AMF0Null,
		AMF0ECMAArray,
		AMF0ObjectEnd,
		AMF0StrictArray,
		AMF0Date,
		AMF0LongString,
		AMF0XMLDocument,
		AMF0TypedObject,
		AMF0SwitchToAMF3}

	for i := 0; i < len(AMF0Types)-1; i++ {
		if AMF0Types[i] == AMF0(b) {
			return AMF0Types[i], true
		}
	}

	return AMF0Null, false
}
