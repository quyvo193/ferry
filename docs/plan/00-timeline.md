# Timeline

~20 h/week, split into 5–6 sessions of 2–4 h. Dates are suggestions. Shift the whole plan if you start later.
Each session lists its **exit criterion**: what must be true before you stop.

---

## Week 1 (Nov 16–22): Foundations: M1 + first half of M2 (20 h)

| # | Hours | Work | Exit criterion |
|---|---|---|---|
| 1.1 | 3 | Push skeleton to GitHub, CI green. Read `net`, `io` docs. M1: receiver accept loop + sender dial/copy | 1 file goes A→B on loopback |
| 1.2 | 3 | M1: filename header, `--dir`, timing + MB/s output, loopback test | Test passes; throughput printed both sides |
| 1.3 | 2 | M1: buffer size sweep, iperf3 baseline on two real machines, fill BENCHMARKS.md | Table filled → **tag `m1`** |
| 1.4 | 4 | M2: frame reader/writer, message types, encode/decode + round-trip tests | `go test ./internal/protocol` green |
| 1.5 | 4 | M2: HELLO handshake, file-list walking (files + folders), OFFER/ACCEPT (auto-accept) | Sender prints offer, receiver prints decoded offer |
| 1.6 | 4 | M2: DATA streaming + SHA-256 on both sides + DONE with hash | Single file over framed protocol, hash verified |

**Week 1 checkpoint:** a single file goes over the framed protocol with hash verification.
If not there: drop the buffer sweep detail (keep 1 MiB) and move on.

---

## Week 2 (Nov 23–29): Finish M2, M3 discovery (20 h)

| # | Hours | Work | Exit criterion |
|---|---|---|---|
| 2.1 | 3 | M2: path sanitization (`os.Root` / `filepath.IsLocal`), `.part` + fsync + rename, name conflict policy | Malicious path test rejected; folders recreated |
| 2.2 | 3 | M2: integration test (folder tree), fuzz frame decoder, PROTOCOL.md v0.1, re-measure | **tag `m2`** |
| 2.3 | 3 | M3 spike: multicast send/recv between 2 real machines (throwaway code) | Packet seen on other machine |
| 2.4 | 4 | M3: announcement packet codec, announcer + listener goroutines, peer table with TTL | Peers appear/expire in logs |
| 2.5 | 3 | M3: `ferry peers`, peer resolution in `ferry send <name>` | Send by name works → **tag `m3`** |
| 2.6 | 4 | **Buffer**: catch up / refactor / CI on 3 OSes. Read about TLS 1.3 + x509 for M4 | Repo clean, notes on TLS written |

**Week 2 checkpoint:** `ferry send <name> <folder>` works on the LAN (no encryption yet).

---

## Week 3 (Nov 30–Dec 6): Security: M4 + start M5 (20 h)

| # | Hours | Work | Exit criterion |
|---|---|---|---|
| 3.1 | 4 | M4: identity (key + self-signed cert generate/persist/load, fingerprint, device ID); wrap conn in TLS | Transfer works over TLS (no pinning yet) |
| 3.2 | 4 | M4: mutual TLS, pinning with `VerifyConnection`, trust store (known_peers) | Unknown/changed peer detected |
| 3.3 | 4 | M4: PIN derivation + pairing UX, accept prompt with timeout, REJECT path | Full first-contact flow works |
| 3.4 | 4 | M4: tests (pinned / mismatch / reject), fingerprint in discovery, TLS throughput, PROTOCOL.md | **tag `m4`** |
| 3.5 | 4 | M5: resume, `.part` detection + metadata, offsets in ACCEPT, prefix re-hash | Kill sender mid-transfer → rerun resumes |

**Week 3 checkpoint:** secure, paired, approved transfers. This is the minimum shippable product.

---

## Week 4 (Dec 7–13): Polish and ship: M5 + release (20 h)

| # | Hours | Work | Exit criterion |
|---|---|---|---|
| 4.1 | 2 | M5: resume edge cases + fault-injection test | Resume test green |
| 4.2 | 4 | M5: progress package (speed EMA, %, ETA, TTY detection) on both sides | Live progress line |
| 4.3 | 4 | M5: measure single stream vs iperf3 → small-file batching → parallel streams (**time-boxed**) | Decision recorded in BENCHMARKS.md |
| 4.4 | 4 | M5: finish parallel streams **or cut** (see cut list); `--streams` flag | **tag `m5`** |
| 4.5 | 3 | Release: GoReleaser snapshot, tag `v0.1.0`, verify `go install` + binaries on 3 OSes | Release page has 6 archives |
| 4.6 | 3 | Benchmarks Ethernet vs Wi-Fi, finalize PROTOCOL.md + README | **v0.1.0 done** |

---

## Weekly ritual (15 min, end of each week)

- Update hours actual vs estimate
- Re-check [09-risks-cut-list.md](09-risks-cut-list.md): cut something now rather than in week 4
- Write 3–5 lines in a dev log (what you learned, what surprised you)
