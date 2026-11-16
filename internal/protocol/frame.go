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
	HeaderSize     = 5
	MaxPayloadSize = 64 << 10
)

type Writer struct {
	bw  *bufio.Writer
	hdr [HeaderSize]byte
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		bw: bufio.NewWriterSize(w, MaxPayloadSize),
	}
}

func (fw *Writer) Write(typ byte, payload []byte) error {
	if len(payload) > MaxPayloadSize {
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

type Reader struct {
	br  *bufio.Reader
	hdr [HeaderSize]byte
	buf []byte
}

func NewReader(r io.Reader) *Reader {
	return &Reader{br: bufio.NewReaderSize(r, MaxPayloadSize)}
}

func (fr *Reader) Read() (typ byte, payload []byte, err error) {
	if _, err := io.ReadFull(fr.br, fr.hdr[:]); err != nil {
		return 0, nil, err
	}

	size := binary.BigEndian.Uint32(fr.hdr[:4])
	if size > MaxPayloadSize {
		return 0, nil, fmt.Errorf("%w: %d bytes", ErrFrameTooLarge, size)
	}

	if cap(fr.buf) < int(size) {
		fr.buf = make([]byte, size)
	}

	payload = fr.buf[:size]
	if _, err := io.ReadFull(fr.br, payload); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return 0, nil, err
	}

	return fr.hdr[4], payload, nil
}
