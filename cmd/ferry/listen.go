package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/quyvo193/ferry/internal/transfer"
)

// runListen makes this device discoverable and receives incoming transfers.
//
// TODO(M1): accept a TCP connection on a fixed port and write the incoming file.
// TODO(M3): start the discovery announcer.
// TODO(M4): wrap the listener in TLS and prompt the user to accept offers.
func runListen(args []string) error {
	fs := flag.NewFlagSet("listen", flag.ContinueOnError)

	port := fs.Int("port", 6767, "TCP port")
	dir := fs.String("dir", ".", "directory to receive files")

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

	if err := os.MkdirAll(*dir, 0o755); err != nil {
		return err
	}

	ln, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(*port)))
	if err != nil {
		return err
	}

	defer ln.Close()
	r := &transfer.Receiver{Dir: *dir}
	return r.Serve(ln)
}
