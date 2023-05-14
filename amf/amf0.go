package amf

const (
	AMF0NumberMarker    = 0x00
	AMF0BooleanMarker   = 0x01
	AMF0StringMarker    = 0x02
	AMF0ObjectMarker    = 0x03
	AMF0NullMarker      = 0x05
	AMF0UndefinedMarker = 0x06
	AMF0ReferenceMarker = 0x07
	AMF0EcmaArrayMarker = 0x08
	AMF0ObjectEndMarker = 0x09
)

type Object map[string]interface{}
