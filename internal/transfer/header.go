package transfer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const maxNameLen = 255

type Header struct {
	Name string
	Size int64
}

var (
	ErrInvalidName = errors.New("invalid file name")
	ErrInvalidSize = errors.New("invalid file size")
)

func readHeader(r io.Reader) (Header, error) {
	var lenBuf [2]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return Header{}, fmt.Errorf("read name length: %w", err)
	}
	lenName := binary.BigEndian.Uint16(lenBuf[:])

	if lenName == 0 || lenName > maxNameLen {
		return Header{}, fmt.Errorf("%w: length %d", ErrInvalidName, lenName)
	}

	nameBuf := make([]byte, lenName)
	if _, err := io.ReadFull(r, nameBuf); err != nil {
		return Header{}, fmt.Errorf("read name: %w", err)
	}

	name := string(nameBuf)
	if !isSafeName(name) {
		return Header{}, fmt.Errorf("%w: %q", ErrInvalidName, name)
	}

	var sizeBuf [8]byte
	if _, err := io.ReadFull(r, sizeBuf[:]); err != nil {
		return Header{}, fmt.Errorf("read size file: %w", err)
	}
	size := binary.BigEndian.Uint64(sizeBuf[:])

	if size > math.MaxInt64 {
		return Header{}, fmt.Errorf("%w: %d", ErrInvalidSize, size)
	}

	return Header{name, int64(size)}, nil
}

func isSafeName(name string) bool {
	return utf8.ValidString(name) &&
		name != "." &&
		filepath.IsLocal(name) &&
		!strings.ContainsAny(name, `/\`)
}

func writeHeader(w io.Writer, h Header) error {
	if len(h.Name) == 0 || len(h.Name) > maxNameLen || !isSafeName(h.Name) {
		return fmt.Errorf("%w: %q", ErrInvalidName, h.Name)
	}

	if h.Size < 0 {
		return fmt.Errorf("%w: %d", ErrInvalidSize, h.Size)
	}

	buf := make([]byte, 0, 2+len(h.Name)+8)
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(h.Name)))
	buf = append(buf, h.Name...)
	buf = binary.BigEndian.AppendUint64(buf, uint64(h.Size))

	_, err := w.Write(buf)

	return err
}
