package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"

	"github.com/quyvo193/ferry/internal/transfer"
)

// runListen makes this device discoverable and receives incoming transfers.
//
// TODO(M3): start the discovery announcer.
// TODO(M4): wrap the listener in TLS and prompt the user to accept offers.
func runListen(args []string) error {
	fs := flag.NewFlagSet("listen", flag.ContinueOnError)

	port := fs.Int("port", 6767, "TCP port")
	dir := fs.String("dir", ".", "directory to receive files")
	discard := fs.Bool("discard", false, "read incoming data and throw it away (network-only benchmark)")
	buf := fs.Int("buf", 0, "copy buffer size in bytes (0 = io.Copy default)")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	if fs.NArg() > 0 {
		return fmt.Errorf("unexpected arguments: %v", fs.Args())
	}

	if *port < 1 || *port > 65535 {
		return fmt.Errorf("invalid port %d", *port)
	}

	if *buf < 0 {
		return fmt.Errorf("invalid buffer size %d", *buf)
	}

	if err := os.MkdirAll(*dir, 0o755); err != nil {
		return err
	}

	ln, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(*port)))
	if err != nil {
		return err
	}

	defer ln.Close()
	slog.Info("listening on", "port", *port, "dir", *dir, "discard", *discard, "buf", *buf)

	r := &transfer.Receiver{Dir: *dir, Discard: *discard, BufSize: *buf}
	return r.Serve(ln)
}
