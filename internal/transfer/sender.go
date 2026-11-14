package transfer

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Sender struct {
	BufSize int
}

func (s *Sender) Send(addr, path string) error {
	f, err := os.Open(path)

	if err != nil {
		return fmt.Errorf("open file failed: %w", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	h := Header{
		Name: fi.Name(),
		Size: fi.Size(),
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}
	defer conn.Close()

	return s.sendFile(conn, f, h)
}

// Private: chỉ Send (và test) gọi.
func (s *Sender) sendFile(conn net.Conn, src io.Reader, h Header) error {
	// writeHeader, CopyN, CloseWrite, đợi EOF, log
	return nil
}
