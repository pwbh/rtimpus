package rtimpus

type Connection struct {
	Handshake HandshakeState
	Err       error
}

func (c *Connection) ProcessMessage(message []byte) {
	switch c.Handshake {
	case Uninitialized:
		handleUninitialized(message)
	}
}

func handleUninitialized(message []byte) {
	result := isVersionSupported(message)
}

func isVersionSupported(message []byte) bool {
	if len(message) > 1 {
		return false
	}

	return SUPPORTED_VERSION == message[0]

}
