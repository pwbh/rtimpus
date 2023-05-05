package rtimpus

import (
	"fmt"
	"io"
	"net"
	"os"
)

const SUPPORTED_PROTOCOL_VERSION = byte(3)

// StartRTMPServer will start an RTMP server on the provided address
func StartRTMPListener(address string) (*net.TCPListener, error) {
	laddr, err := net.ResolveTCPAddr("tcp", address)

	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", laddr)

	if err != nil {
		return nil, err
	}

	return listener, nil
}

func LoopConnections(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			fmt.Fprintf(os.Stderr, "something wen't wrong couldn't accept TCP connection:\n%v", err)
			continue
		}

		go listenOnConnection(conn)
	}
}

func Close(listener *net.TCPListener) {
	listener.Close()
}

func listenOnConnection(tcpConn *net.TCPConn) *Connection {
	connection := new(Connection)
	connection.conn = tcpConn

	go func() {
		buf := make([]byte, 0, 1024)

		for {
			n, err := tcpConn.Read(buf)

			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Fprintf(os.Stderr, "something wen't wrong during exchange on information:\n%v", err)
				connection.Err = err
				return
			}

			fmt.Println(buf[:n])
		}
	}()

	return connection
}
