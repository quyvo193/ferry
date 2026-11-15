package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/quyvo193/ferry/internal/transfer"
)

// runSend sends one or more files/folders to a peer.
//
// TODO(M2): build an OFFER from files/folders and stream them with hashing.
// TODO(M3): resolve <peer> by name or device ID via discovery.
func runSend(args []string) error {
	fs := flag.NewFlagSet("send", flag.ContinueOnError)
	buf := fs.Int("buf", 0, "copy buffer size in bytes (0 = io.Copy default, uses sendfile when possible)")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if fs.NArg() != 2 {
		return fmt.Errorf("usage: ferry send [-buf N] <ip:port> <file>")
	}
	if *buf < 0 {
		return fmt.Errorf("invalid buffer size %d", *buf)
	}

	s := &transfer.Sender{BufSize: *buf}
	addr, path := fs.Arg(0), fs.Arg(1)
	return s.Send(addr, path)
}
