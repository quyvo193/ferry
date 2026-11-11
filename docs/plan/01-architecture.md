# Architecture

This is a starting point, not a contract. Change it when the code tells you to, and update this file when you do.

## Package layout and dependencies

```
cmd/ferry  ──►  transfer  ──►  protocol
    │              │  └────►  progress
    │              └───────►  trust ──► identity ──► config
    ├────────►  discovery ─────────────► identity
    └────────►  config
```

Rules:
- `protocol` depends only on the standard library. It knows bytes and messages, nothing about files or users.
- `transfer` owns file-system logic (walking, `.part`, rename, resume) and drives `protocol`.
- `cmd/ferry` is the only place that touches stdin/stdout prompts and flags. Lower packages take
  interfaces/callbacks (e.g. `Approver` for the accept prompt) so they stay testable.
- Everything under `internal/` so nobody imports it before the API is stable.

| Package | Responsibility | Milestone |
|---|---|---|
| `config` | Config dir, settings (name, port, download dir) | M3–M4 |
| `identity` | Key + self-signed cert, fingerprint, device ID | M4 |
| `trust` | known_peers store, PIN derivation, TLS verify hook | M4 |
| `discovery` | Multicast announcer/listener, peer table with TTL | M3 |
| `protocol` | Frame codec, message encode/decode | M2 |
| `transfer` | Send/receive sessions, hashing, safe writes, resume, parallel | M1–M5 |
| `progress` | Byte counters, speed/ETA rendering | M5 |

## Connection lifecycle (target state after M4)

```
sender                                        receiver
──────                                        ────────
resolve peer (discovery table / ip:port)
TCP dial ───────────────────────────────────► accept
TLS 1.3 handshake (mutual, self-signed) ◄───► TLS handshake
  VerifyConnection: pinned? / pairing          VerifyConnection: pinned? / pairing
HELLO(version, device_id, name) ────────────► check version, device_id ↔ cert
                              ◄──────────────  HELLO
[first contact: both show PIN, user confirms]
OFFER(files, sizes) ────────────────────────► Approver: prompt user (timeout)
                    ◄────────────────────────  ACCEPT(offsets) | REJECT(reason)
DATA… DONE(file, hash) … DONE(all) ─────────► write .part, hash, verify, rename
                    ◄────────────────────────  final result / ERROR
close                                          close
```

## Concurrency model

- **Receiver (`listen`)**: one accept loop goroutine, one goroutine per connection. The accept
  prompt is serialized (one prompt at a time) because stdin is shared. Decide: queue or reply "busy".
- **Discovery**: announcer goroutine (ticker) + listener goroutine + table sweeper. All stop via a
  shared `context.Context`.
- **Transfer**: one stream = one goroutine reading/writing sequentially. Parallel streams (M5) =
  N goroutines with an `errgroup`-style "first error cancels all" pattern. Stdlib only:
  `sync.WaitGroup` + `context.WithCancelCause`.
- **Shutdown**: `signal.NotifyContext(ctx, os.Interrupt)` in `main`, passed down everywhere.
  Ctrl-C must leave `.part` files (resumable), never half-renamed files.

## On-disk layout

```
<os.UserConfigDir()>/ferry/        # macOS: ~/Library/Application Support/ferry
  config.json                      #   name, port, download dir (optional, flags override)
  identity/
    cert.pem                       #   self-signed cert
    key.pem                        #   private key (0600)
  known_peers.json                 #   device_id → fingerprint, name, first_seen

<download dir>/                    # default: ~/Downloads/ferry (decision: or current dir?)
  photos/a.jpg
  photos/b.jpg.part                #   in progress
  photos/b.jpg.part.meta           #   resume metadata (M5, if you choose sidecar)
```

## Suggested defaults (change if you want)

| Setting | Value | Notes |
|---|---|---|
| Discovery group | `239.255.70.70` (admin-scoped) | Avoid LocalSend's `224.0.0.167` to not collide |
| Discovery port | UDP `47700` | |
| Transfer port | TCP `47701` | `--port` flag |
| Announce interval | 2 s | + small random jitter |
| Peer TTL | 7 s | ≈ 3 missed announcements |
| Copy buffer | 1 MiB | tune in M1 |
| Max frame size | 16 MiB | DoS guard |
| Hash | SHA-256 | hardware-accelerated on amd64/arm64 |

## Conventions

- **Errors**: wrap with context `fmt.Errorf("send offer: %w", err)`. Sentinel errors for protocol
  conditions (`ErrFrameTooLarge`, `ErrVersionMismatch`) so tests can `errors.Is`.
- **Logging**: `log/slog` to stderr; `--verbose` switches to debug level. User-facing output
  (prompts, progress) goes to stdout and is separate from logs.
- **Sizes**: `int64`/`uint64` for file sizes and offsets everywhere. Files > 4 GiB are a test case.
- **Flags**: one `flag.FlagSet` per subcommand.
- **Commits**: `feat(protocol): …`, `fix(transfer): …`, `docs: …`. Works with the GoReleaser changelog filter.
