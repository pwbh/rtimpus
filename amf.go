package rtimpus

import "fmt"

type AMF0 byte

const (
	AMF0Number       = 0x00 // Number
	AMF0Boolean      = 0x01 // Boolean
	AMF0Stirng       = 0x02 // String
	AMF0Object       = 0x03 // Object (Set of key/value pairs)
	AMF0Null         = 0x05 // Null
	AMF0ECMAArray    = 0x08 // ECMA Array
	AMF0ObjectEnd    = 0x09 // Object End
	AMF0StrictArray  = 0x0a // Strict Array
	AMF0Date         = 0x0b // Date
	AMF0LongString   = 0x0c // Long String
	AMF0XMLDocument  = 0x0f // XML Document
	AMF0TypedObject  = 0x10 // Typed Object
	AMF0SwitchToAMF3 = 0x11 // Switch to AMF3
)

type Value struct {
	Type AMF0
	Data []byte
}

type AMF0Result struct {
	Command       string
	TransactionID uint32
	Payload       []Value
}

func UnmarshalAMF0(message []byte) *AMF0Result {
	result := new(AMF0Result)

	phase := 0
	start := 0
	end := 0

	for len(message) > end {
		offset, ok := getFoundOffset(message, end)

		if ok && end > start {
			if phase == 0 {
				fmt.Println(message[start:end])
				phase = 1
			} else if phase == 1 {
				fmt.Println(message[start:end])
				phase = 2
			} else {
				fmt.Println(message[start:end])
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

func getFoundOffset(message []byte, currentIndex int) (int, bool) {
	a := isAmfType(message[currentIndex])

	if !a {
		return 0, false
	}

	b := false

	if len(message) > currentIndex+1 {
		b = message[currentIndex+1] == byte(AMF0Number)
	}

	if a && b {
		return 2, true
	}

	return 1, true
}

func isAmfType(b byte) bool {
	AMF0Types := [...]byte{
		AMF0Number,
		AMF0Boolean,
		AMF0Stirng,
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
		if AMF0Types[i] == b {
			return true
		}
	}

	return false
}
