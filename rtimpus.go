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
	connection.Phase = Uninitialized

	go func() {
		buf := make([]byte, 0, 4096)
		tmp := make([]byte, 256)

		total := 0

		for {

			n, err := tcpConn.Read(tmp)

			if err != nil {
				if err == io.EOF {
					fmt.Fprintf(os.Stdout, "received IOF:\n%v\n", err)
					break
				}
				fmt.Fprintf(os.Stderr, "something wen't wrong during exchange of information:\n%v\n", err)
				connection.Err = err
				return
			}

			total += n
			buf = append(buf, tmp[:n]...)

			if len(tmp) > n {
				connection.Process(buf[:total])
				total = 0
				buf = make([]byte, 0, 4096)
			}
		}
	}()

	return connection
}
