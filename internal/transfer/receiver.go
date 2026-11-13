package transfer

import (
	"errors"
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
			if errors.Is(err, net.ErrClosed) {
				return nil
			}

			log.Printf("accept connection fail %v\n", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			if err := r.receive(conn); err != nil {
				log.Printf("%s: %v", conn.RemoteAddr(), err)
			}
		}(conn)
	}
}

func (r *Receiver) receive(src io.Reader) error {
	fname := make([]byte, 1024)
	fsize := make([]byte, 1024)
	_, err := io.ReadFull(src, fname)
	_, err = io.ReadFull(src, fsize)

	if err != nil {
		return err
	}

	return nil
}
