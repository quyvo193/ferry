# ferry

Send files and folders directly between devices on the same Wi-Fi/LAN. No cloud, no account.

Devices find each other automatically and transfer over an encrypted peer-to-peer connection.
The goal is speed (saturate the local network) and learning how it all works underneath.

> **Status: early development.** Nothing works yet. See the [roadmap](#roadmap).

## Features (planned for v1)

- **Zero-config discovery**: devices announce themselves via UDP multicast
- **Encrypted**: TLS 1.3 everywhere, self-signed per-device certificates, trust-on-first-use pairing with a short PIN
- **Receiver approval**: you see who's sending what (names + sizes) before any byte is transferred
- **Fast**: custom binary protocol, large buffers, no per-chunk ACKs, optional parallel streams
- **Folders**: multiple files and whole directory trees, structure preserved
- **Integrity**: every file is hashed on both sides and verified before it's committed
- **Safe writes**: data lands in `file.part` and is renamed only after verification
- **Resume**: interrupted transfers continue from where they stopped
- **Progress**: live speed, percentage and ETA
- **Single static binary** for macOS, Windows and Linux

## Install

### With Go

Requires Go 1.27+.

```sh
go install github.com/quyvo193/ferry/cmd/ferry@latest
```

The binary goes to `$(go env GOPATH)/bin`. Make sure that directory is on your `PATH`.

### Prebuilt binaries

Download the archive for your OS/arch from the [Releases](https://github.com/quyvo193/ferry/releases) page,
extract it, and put `ferry` somewhere on your `PATH`.

The binaries aren't signed:

- **macOS**: `xattr -d com.apple.quarantine ./ferry`
- **Windows**: SmartScreen may warn you. Choose *More info → Run anyway*.

### From source

```sh
git clone https://github.com/quyvo193/ferry.git
cd ferry
make build        # -> ./bin/ferry
make install      # -> $(go env GOPATH)/bin/ferry
```

## Usage

```sh
# On the receiving device: become discoverable and wait for transfers
ferry listen

# On the sending device: see who's around
ferry peers

# Send files and folders to a peer (by name or device ID)
ferry send "Quy's MacBook" ./photos ./notes.txt
```

| Command | Description |
|---|---|
| `ferry listen` | Become discoverable and receive files |
| `ferry peers` | List discovered devices |
| `ferry send <peer> <path>...` | Send files/folders to a peer |
| `ferry version` | Print version |

The first time two devices talk, both show a short PIN. Confirm it matches and the devices
remember each other.

> Your OS firewall may ask to allow incoming connections the first time you run `ferry listen`.

## How it works

1. **Discovery**: each device multicasts a small UDP announcement every few seconds
   (device ID, name, TCP port, certificate fingerprint). Peers that go quiet expire after a TTL.
2. **Identity and trust**: on first run each device generates a self-signed TLS certificate.
   The first connection between two devices is confirmed with a PIN derived from both
   certificate fingerprints. After that the fingerprint is pinned (trust on first use, like SSH).
3. **Protocol**: length-prefixed binary frames over TLS/TCP:
   `[u32 length][u8 type][payload]`, big-endian.
   Messages: `HELLO`, `OFFER`, `ACCEPT`/`REJECT`, `DATA`, `DONE`, `ERROR`.
   Full spec: [docs/PROTOCOL.md](docs/PROTOCOL.md).
4. **Transfer**: raw bytes are streamed with large buffers. Both sides hash while streaming
   and compare at the end. Files are written to `.part` and renamed after verification.

## Project layout

```
cmd/ferry/            CLI entry point and subcommands
internal/config/      paths and settings
internal/identity/    device certificate, key, fingerprint, device ID
internal/trust/       known peers, TOFU pairing, PIN
internal/discovery/   UDP multicast announce + peer table
internal/protocol/    frame codec and message types
internal/transfer/    send/receive sessions, hashing, safe writes, resume
internal/progress/    speed / percentage / ETA rendering
docs/                 protocol spec, benchmarks, development plan
```

## Development

```sh
make help     # list targets
make check    # gofmt + vet + tests (same as CI)
make race     # tests with the race detector
make cross    # cross-compile all platforms into ./dist
```

Releases: push a tag `vX.Y.Z` and GitHub Actions runs [GoReleaser](https://goreleaser.com).

## Roadmap

- [ ] M1: raw TCP send/receive with a manual IP, throughput baseline
- [ ] M2: framing protocol, multiple files and folders, hashing, safe writes
- [ ] M3: UDP multicast discovery
- [ ] M4: TLS, pairing, accept prompt
- [ ] M5: resume, parallel streams, progress
- [ ] v0.1.0 release + benchmarks (Ethernet vs Wi-Fi)

Detailed plan: [docs/plan](docs/plan/README.md). Benchmarks: [docs/BENCHMARKS.md](docs/BENCHMARKS.md).

## Out of scope (for now)

Bluetooth, mobile apps, transfers over the internet/relay, GUI.

## About

A learning project. All code is written by hand; AI is only used to explain concepts and review code.

## License

[MIT](LICENSE)
