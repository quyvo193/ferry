package transfer

import (
	"io"
	"log/slog"
	"time"
)

func copyN(dst io.Writer, src io.Reader, n int64, bufSize int) (int64, error) {
	if bufSize <= 0 {
		return io.CopyN(dst, src, n)
	}

	written, err := io.CopyBuffer(
		struct{ io.Writer }{dst},
		struct{ io.Reader }{io.LimitReader(src, n)},
		make([]byte, bufSize),
	)
	if err == nil && written < n {
		err = io.EOF
	}
	return written, err
}

func logRate(msg, name string, n int64, d time.Duration) {
	secs := d.Seconds()
	slog.Info(msg,
		"name", name,
		"bytes", n,
		"secs", secs,
		"MB/s", float64(n)/1e6/secs,
		"Mbit/s", float64(n)*8/1e6/secs,
	)
}
