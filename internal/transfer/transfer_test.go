package transfer

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	os.Exit(m.Run())
}

func startReceiver(t *testing.T, r *Receiver) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- r.Serve(ln) }()
	t.Cleanup(func() {
		ln.Close()
		if err := <-done; err != nil {
			t.Errorf("Serve: %v", err)
		}
	})
	return ln.Addr().String()
}

func writeRandom(t *testing.T, path string, size int) []byte {
	t.Helper()
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return data
}

func TestLoopback(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		size    int
		bufSize int
	}{
		{"10MiB default copy", "random.bin", 10 << 20, 0},
		{"10MiB 32KiB buffer", "random.bin", 10 << 20, 32 << 10},
		{"10MiB 1MiB buffer", "random.bin", 10 << 20, 1 << 20},
		{"buffer larger than file", "small.bin", 1000, 1 << 20},
		{"zero bytes", "empty.txt", 0, 0},
		{"zero bytes with buffer", "empty.txt", 0, 32 << 10},
		{"spaces", "my file (1).txt", 4096, 0},
		{"unicode", "ảnh chụp 日本語.png", 4096, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcDir, dstDir := t.TempDir(), t.TempDir()
			src := filepath.Join(srcDir, tt.file)
			data := writeRandom(t, src, tt.size)

			addr := startReceiver(t, &Receiver{Dir: dstDir, BufSize: tt.bufSize})
			s := &Sender{BufSize: tt.bufSize}
			if err := s.Send(addr, src); err != nil {
				t.Fatalf("Send: %v", err)
			}

			got, err := os.ReadFile(filepath.Join(dstDir, tt.file))
			if err != nil {
				t.Fatal(err)
			}
			if sha256.Sum256(got) != sha256.Sum256(data) {
				t.Fatalf("hash mismatch: got %d bytes, want %d", len(got), len(data))
			}

			entries, err := os.ReadDir(dstDir)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 1 {
				t.Fatalf("dst dir has %d entries, want 1", len(entries))
			}
		})
	}
}

func TestLoopbackDiscard(t *testing.T) {
	srcDir, dstDir := t.TempDir(), t.TempDir()
	src := filepath.Join(srcDir, "random.bin")
	writeRandom(t, src, 1<<20)

	addr := startReceiver(t, &Receiver{Dir: dstDir, Discard: true})
	if err := (&Sender{}).Send(addr, src); err != nil {
		t.Fatalf("Send: %v", err)
	}

	entries, err := os.ReadDir(dstDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("discard wrote %d entries", len(entries))
	}
}

func TestLoopbackDoesNotOverwrite(t *testing.T) {
	srcDir, dstDir := t.TempDir(), t.TempDir()
	src := filepath.Join(srcDir, "a.txt")
	writeRandom(t, src, 4096)

	existing := []byte("keep me")
	dst := filepath.Join(dstDir, "a.txt")
	if err := os.WriteFile(dst, existing, 0o644); err != nil {
		t.Fatal(err)
	}

	addr := startReceiver(t, &Receiver{Dir: dstDir})
	_ = (&Sender{}).Send(addr, src)

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, existing) {
		t.Fatal("existing file was overwritten")
	}
	if _, err := os.Stat(dst + ".part"); !os.IsNotExist(err) {
		t.Fatalf(".part left behind: %v", err)
	}
}

func TestReceiveTruncatedRemovesPart(t *testing.T) {
	dstDir := t.TempDir()
	h := Header{Name: "cut.bin", Size: 1000}

	var buf bytes.Buffer
	if err := writeHeader(&buf, h); err != nil {
		t.Fatal(err)
	}
	buf.Write(make([]byte, 400))

	r := &Receiver{Dir: dstDir}
	if err := r.receive(&buf); err == nil {
		t.Fatal("receive accepted a truncated stream")
	}

	entries, err := os.ReadDir(dstDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("dst dir has %d entries after failed receive, want 0", len(entries))
	}
}

func TestSendRejectsDirectory(t *testing.T) {
	if err := (&Sender{}).Send("127.0.0.1:1", t.TempDir()); err == nil {
		t.Fatal("Send accepted a directory")
	}
}
