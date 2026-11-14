package transfer

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"
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

			slog.Error("accept connection failed", "err", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			if err := r.receive(conn); err != nil {
				slog.Error("receive failed", "remote", conn.RemoteAddr().String(), "err", err)
			}
		}(conn)
	}
}

func (r *Receiver) receive(src io.Reader) error {
	h, err := readHeader(src)

	if err != nil {
		return err
	}

	final := filepath.Join(r.Dir, h.Name)
	if _, err := os.Lstat(final); err == nil {
		return fmt.Errorf("%s: already exists", final)
	}

	f, err := os.Create(h.Name)
	if err != nil {
		return err
	}

	part := final + ".part"
	done := false
	defer func() {
		if !done {
			f.Close()
			os.Remove(part)
		}
	}()

	start := time.Now()

	n, err := io.CopyN(f, src, h.Size)
	if err != nil {
		return fmt.Errorf("copy data: %w", err)
	}

	if err := os.Rename(part, final); err != nil {
		return err
	}
	done = true
	secs := time.Since(start).Seconds()
	slog.Info("received",
		"name", h.Name,
		"bytes", n,
		"secs", secs,
		"MB/s", float64(n)/1e6/secs,
	)

	return nil
}
