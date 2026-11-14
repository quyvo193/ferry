package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/quyvo193/ferry/internal/transfer"
)

// runSend sends one or more files/folders to a peer.
//
// TODO(M1): dial a manually entered ip:port and stream a single file.
// TODO(M2): build an OFFER from files/folders and stream them with hashing.
// TODO(M3): resolve <peer> by name or device ID via discovery.
func runSend(args []string) error {
	fs := flag.NewFlagSet("send", flag.ContinueOnError)

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	addr, path := fs.Arg(0), fs.Arg(1)

	if fs.NArg() != 2 {
		return fmt.Errorf("usage: ferry send <ip:port> <file>")
	}

	s := &transfer.Sender{BufSize: 32 << 10}
	return s.Send(addr, path)
}
