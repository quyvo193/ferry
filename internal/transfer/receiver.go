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
		return fmt.Errorf("read header: %w", err)
	}

	start := time.Now()
	n, err := r.store(src, h)
	if err != nil {
		return err
	}

	secs := time.Since(start).Seconds()
	slog.Info("received",
		"name", h.Name,
		"bytes", n,
		"secs", secs,
		"MB/s", float64(n)/1e6/secs,
	)

	return nil
}

func (r *Receiver) store(src io.Reader, h Header) (int64, error) {
	if r.Discard {
		return io.CopyN(io.Discard, src, h.Size)
	}

	final := filepath.Join(r.Dir, h.Name)
	if _, err := os.Lstat(final); err == nil {
		return 0, fmt.Errorf("%s: already exists", final)
	}

	part := final + ".part"
	f, err := os.OpenFile(part, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return 0, err
	}
	done := false
	defer func() {
		if !done {
			f.Close()
			os.Remove(part)
		}
	}()

	n, err := io.CopyN(f, src, h.Size)
	if err != nil {
		return n, fmt.Errorf("copy data: got %d/%d bytes: %w", n, h.Size, err)
	}

	if err := f.Close(); err != nil {
		return n, fmt.Errorf("close %s: %w", part, err)
	}
	if err := os.Rename(part, final); err != nil {
		return n, err
	}
	done = true
	return n, nil
}
