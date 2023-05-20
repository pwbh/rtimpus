package rtimpus

import (
	"fmt"
	"io"
	"net"
	"os"
)

const SUPPORTED_PROTOCOL_VERSION = byte(3)

// StartRTMPServer will start an RTMP server on the provided address
func Listen(address string) (*Listener, error) {
	laddr, err := net.ResolveTCPAddr("tcp", address)

	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", laddr)

	if err != nil {
		return nil, err
	}

	rtmpListner := &Listener{tcpListener: listener}

	return rtmpListner, nil
}

func LoopConnections(listener *Listener) {
	for {
		conn, err := listener.tcpListener.AcceptTCP()

		if err != nil {
			fmt.Fprintf(os.Stderr, "something wen't wrong couldn't accept TCP connection:\n%v", err)
			continue
		}

		go listenOnConnection(conn)
	}
}

func listenOnConnection(tcpConn *net.TCPConn) *Connection {
	connection := new(Connection)
	connection.Conn = tcpConn
	connection.Phase = Uninitialized
	connection.ChunkSize = 128

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
