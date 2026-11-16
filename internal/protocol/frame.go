package protocol

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	ErrFrameTooLarge = errors.New("protocol: frame too large")
)

const (
	HeaderSize   = 5
	MaxFrameSize = 1<<20 + 64
)

type Writer struct {
	bw  *bufio.Writer
	hdr [HeaderSize]byte
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		bw: bufio.NewWriterSize(w, MaxFrameSize),
	}
}

func (fw *Writer) Write(typ byte, payload []byte) error {
	if len(payload) > MaxFrameSize {
		return fmt.Errorf("%w: %d bytes", ErrFrameTooLarge, len(payload))
	}

	// write payload size to header
	binary.BigEndian.PutUint32(fw.hdr[0:4], uint32(len(payload)))
	fw.hdr[4] = typ

	// write header
	if _, err := fw.bw.Write(fw.hdr[:]); err != nil {
		return err
	}

	// write payload
	_, err := fw.bw.Write(payload)
	return err
}

func (fw *Writer) Flush() error { return fw.bw.Flush() }
