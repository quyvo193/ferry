package transfer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"testing"
)

func TestHeaderRoundTrip(t *testing.T) {
	tests := []Header{
		{Name: "a.txt", Size: 0},
		{Name: "big.bin", Size: 4 << 30},
		{Name: "my file (1).txt", Size: 42},
		{Name: "ảnh-chụp-màn-hình.png", Size: 1},
		{Name: strings.Repeat("x", maxNameLen), Size: 7},
	}

	for _, want := range tests {
		t.Run(want.Name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := writeHeader(&buf, want); err != nil {
				t.Fatalf("writeHeader: %v", err)
			}
			got, err := readHeader(&buf)
			if err != nil {
				t.Fatalf("readHeader: %v", err)
			}
			if got != want {
				t.Fatalf("got %+v, want %+v", got, want)
			}
			if buf.Len() != 0 {
				t.Fatalf("%d bytes left after header", buf.Len())
			}
		})
	}
}

func rawHeader(name string, size uint64) []byte {
	b := binary.BigEndian.AppendUint16(nil, uint16(len(name)))
	b = append(b, name...)
	return binary.BigEndian.AppendUint64(b, size)
}

func TestReadHeaderRejects(t *testing.T) {
	tests := []struct {
		name string
		raw  []byte
		want error
	}{
		{"empty name", rawHeader("", 1), ErrInvalidName},
		{"name too long", rawHeader(strings.Repeat("x", maxNameLen+1), 1), ErrInvalidName},
		{"dot", rawHeader(".", 1), ErrInvalidName},
		{"dotdot", rawHeader("..", 1), ErrInvalidName},
		{"slash", rawHeader("a/b", 1), ErrInvalidName},
		{"backslash", rawHeader(`a\b`, 1), ErrInvalidName},
		{"absolute", rawHeader("/etc/passwd", 1), ErrInvalidName},
		{"traversal", rawHeader("../x", 1), ErrInvalidName},
		{"invalid utf8", rawHeader("\xff\xfe", 1), ErrInvalidName},
		{"negative size", rawHeader("a", 1<<63), ErrInvalidSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := readHeader(bytes.NewReader(tt.raw))
			if !errors.Is(err, tt.want) {
				t.Fatalf("got %v, want %v", err, tt.want)
			}
		})
	}
}

func TestReadHeaderTruncated(t *testing.T) {
	raw := rawHeader("file.txt", 123)
	for i := range len(raw) {
		if _, err := readHeader(bytes.NewReader(raw[:i])); err == nil {
			t.Fatalf("readHeader accepted %d/%d bytes", i, len(raw))
		}
	}
}
