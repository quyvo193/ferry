# Testing Strategy

Testing time is already inside each milestone's budget. It isn't a separate phase.

## Layers

| Layer | What | Tools | When |
|---|---|---|---|
| Unit | Codecs, path sanitizer, PIN derivation, peer table, progress math | `testing`, table-driven | Every milestone |
| Fuzz | Anything parsing bytes from the network | `go test -fuzz` | M2, M3 |
| Integration | Two in-process nodes over loopback, separate temp dirs | `t.TempDir()`, `net.Listen("tcp", "127.0.0.1:0")` | M1+ |
| Fault injection | Conn that drops/corrupts/slows | custom `net.Conn` wrapper | M2, M5 |
| Race | Concurrency bugs | `go test -race` (CI on Linux/macOS) | Always |
| Benchmarks | Codec, hashing, copy loop | `testing.B`, `-benchmem` | M2, M5 |
| Manual | Real machines, real networks, real OSes | checklist below | End of each milestone |

## Fuzz targets (minimum)

- `FuzzReadFrame`: arbitrary bytes into the frame reader → never panic, never allocate > max frame
- `FuzzDecodeOffer`: arbitrary payload → error or valid offer whose paths pass sanitization
- `FuzzDecodeAnnouncement`: discovery packets

Add any crasher to `testdata/fuzz/…` (Go does this automatically) and commit it.

## Test helpers worth writing once

- **`faultConn`**: wraps `net.Conn`; options: close after N bytes, flip bit at offset K, delay per
  write. Used for corruption (M2) and resume (M5) tests.
- **`testNode`**: spins up a receiver with its own temp config dir + download dir on a random port;
  returns address + cleanup. Makes integration tests one-liners.
- **Tree compare**: walk two dirs, compare relative paths + sizes + SHA-256.
- **Injectable clock** for the peer table TTL (`func() time.Time` field). No `time.Sleep` in tests.

## Edge-case checklist

Files and paths:
- [ ] Empty file, empty folder, deeply nested folder
- [ ] File > 4 GiB (u32 overflow guard)
- [ ] Unicode names (Vietnamese diacritics, emoji), spaces, leading dots
- [ ] Names invalid on Windows when sent from macOS/Linux (`a:b.txt`, `CON`, trailing dot/space)
- [ ] Same filename in different folders; final file already exists (conflict policy)
- [ ] Symlinks (skipped with warning)
- [ ] File deleted or grows during send

Network / session:
- [ ] Sender disconnects mid-file → `.part` only
- [ ] Receiver disconnects mid-file → sender errors out cleanly
- [ ] Wrong protocol version → ERROR, both exit cleanly
- [ ] Oversized frame → rejected
- [ ] Offer rejected / prompt timeout
- [ ] Identity changed → refused

## Manual test matrix (before `v0.1.0`)

| Sender → Receiver | Ethernet | Wi-Fi |
|---|---|---|
| macOS → Linux/Windows | | |
| Linux/Windows → macOS | | |
| same machine, 2 processes | n/a (loopback) | |

For each: discovery visible · pairing PIN · accept prompt · folder transfer verified · resume after kill · progress shown.

## CI

- `ci.yml`: gofmt, vet, tests on ubuntu/macos/windows (race on non-Windows), cross-compile job
- Multicast tests: CI runners often don't allow multicast. Skip with `testing.Short()` or an env
  var (`FERRY_MULTICAST_TESTS=1`) and run them locally.
