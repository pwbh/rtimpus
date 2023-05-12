package rtimpus

import "fmt"

type AMF0 byte

const (
	AMF0Number       = AMF0(0x00)
	AMF0Boolean      = AMF0(0x01)
	AMF0String       = AMF0(0x02)
	AMF0Object       = AMF0(0x03)
	AMF0Null         = AMF0(0x05)
	AMF0ECMAArray    = AMF0(0x08)
	AMF0ObjectEnd    = AMF0(0x09)
	AMF0StrictArray  = AMF0(0x0a)
	AMF0Date         = AMF0(0x0b)
	AMF0LongString   = AMF0(0x0c)
	AMF0XMLDocument  = AMF0(0x0f)
	AMF0TypedObject  = AMF0(0x10)
	AMF0SwitchToAMF3 = AMF0(0x11)
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
	amf0Types :=
		[]byte{0x00,
			0x01,
			0x02,
			0x03,
			0x05,
			0x08,
			0x09,
			0x0a,
			0x0b,
			0x0c,
			0x0f,
			0x10,
			0x11}

	for i := 0; i < len(amf0Types)-1; i++ {
		if amf0Types[i] == b {
			return true
		}
	}

	return false
}
