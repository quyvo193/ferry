// Command ferry sends files and folders directly between devices on the same LAN.
package main

import (
	"errors"
	"fmt"
	"os"
)

const usage = `ferry - send files and folders to devices on your local network.

Usage:
  ferry <command> [flags] [args]

Commands:
  listen                   Become discoverable and receive files
  peers                    List discovered devices
  send <peer> <path>...    Send files or folders to a peer
  version                  Print version information
  help                     Show this help

Run "ferry <command> -h" for command-specific flags.
`

var errNotImplemented = errors.New("not implemented yet")

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	cmd, args := os.Args[1], os.Args[2:]

	var err error
	switch cmd {
	case "listen":
		err = runListen(args)
	case "peers":
		err = runPeers(args)
	case "send":
		err = runSend(args)
	case "version", "-v", "--version":
		printVersion()
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "ferry: unknown command %q\n\n%s", cmd, usage)
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ferry %s: %v\n", cmd, err)
		os.Exit(1)
	}
}
