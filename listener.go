package rtimpus

import "net"

type Listener struct {
	tcpListener *net.TCPListener
}

func (l *Listener) Addr() net.Addr {
	return l.tcpListener.Addr()
}

func (l *Listener) Close() {
	l.tcpListener.Close()
}
