//	Version defined by this specification is 3 (This is the only version this protocol will
//	support 0-2 are depercated values sused by earleir proprietary products,
//	4-31 are reserved for future implementations, 32-35 are not allowed.)

package rtimpus

type Phase byte

const (
	Uninitialized Phase = iota
	VersionSent
	AckSent
	HandshakeDone
)
