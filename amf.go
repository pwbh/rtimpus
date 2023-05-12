package rtimpus

import "fmt"

type AMF0 byte

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
				fmt.Println(string(message[start:end]))
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
		0x00, // Number
		0x01, // Boolean
		0x02, // String
		0x03, // Object (Set of key/value pairs)
		0x05, // Null
		0x08, // ECMA Array
		0x09, // Object End
		0x0a, // Strict Array
		0x0b, // Date
		0x0c, // Long String
		0x0f, // XML Document
		0x10, // Typed Object
		0x11} // Switch to AMF3

	for i := 0; i < len(AMF0Types)-1; i++ {
		if AMF0Types[i] == b {
			return true
		}
	}

	return false
}
