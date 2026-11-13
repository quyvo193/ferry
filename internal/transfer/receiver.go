package transfer

import (
	"io"
	"log"
	"net"
)

type Receiver struct {
	Dir     string
	Discard bool
}

func (r *Receiver) Serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept connection fail %v\n", err)
			continue
		}

		conn.Close()
		go func(conn net.Conn) {
			if err := r.receive(conn); err != nil {
				log.Printf("%s: %v", conn.RemoteAddr(), err)
			}
		}(conn)
	}
}

func (r *Receiver) receive(src io.Reader) error {
	return nil
}
