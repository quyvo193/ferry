package transfer

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"
)

type Sender struct {
	BufSize int
}

func (s *Sender) Send(addr, path string) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if !fi.Mode().IsRegular() {
		return fmt.Errorf("%s: not a regular file", path)
	}

	h := Header{
		Name: fi.Name(),
		Size: fi.Size(),
	}

	if err := validate(h.Name, int(h.Size)); err != nil {
		return err
	}

	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	return s.sendFile(conn, f, h)
}

func (s *Sender) sendFile(conn net.Conn, src io.Reader, h Header) error {
	start := time.Now()
	if err := writeHeader(conn, h); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	n, err := io.CopyN(conn, src, h.Size)
	if err != nil {
		return fmt.Errorf("send data: sent %d/%d bytes: %w", n, h.Size, err)
	}

	if cw, ok := conn.(interface{ CloseWrite() error }); ok {
		if err := cw.CloseWrite(); err != nil {
			return fmt.Errorf("close write: %w", err)
		}
	}

	if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		return err
	}

	if _, err := io.Copy(io.Discard, conn); err != nil {
		return fmt.Errorf("wait for receiver: %w", err)
	}

	secs := time.Since(start).Seconds()
	slog.Info("sent",
		"name", h.Name,
		"bytes", n,
		"secs", secs,
		"MB/s", float64(n)/1e6/secs,
	)

	return nil
}
