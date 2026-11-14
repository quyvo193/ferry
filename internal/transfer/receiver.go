package transfer

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
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

	var dst io.Writer = io.Discard

	start := time.Now()

	n, err := io.CopyN(dst, src, h.Size)
	if err != nil {
		return fmt.Errorf("copy data: %w", err)
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
