package rtimpus

type AMF0 byte

const (
	Number       = AMF0(0x00)
	Boolean      = AMF0(0x01)
	String       = AMF0(0x02)
	Object       = AMF0(0x03)
	Null         = AMF0(0x05)
	ECMAArray    = AMF0(0x08)
	ObjectEnd    = AMF0(0x09)
	StrictArray  = AMF0(0x0a)
	Date         = AMF0(0x0b)
	LongString   = AMF0(0x0c)
	XMLDocument  = AMF0(0x0f)
	TypedObject  = AMF0(0x10)
	SwitchToAMF3 = AMF0(0x11)
)

type Value struct {
	Type AMF0
	Data []byte
}

type AMF0Result struct {
	Command       string
	TransactionID uint32
	Payload       map[string]Value
}

func UnmarshalAMF0(message []byte) *AMF0Result {
	start := 0
	end := 0

	for end < len(message) {
		if isAmfType(message[end]) {
			end += 1
		} else {

		}
	}
}

func isAmfType(b byte) bool {
	return (b > 0x00 && b <= 0x11) && b != 0x04 && b != 0x0d && b != 0x0e
}
