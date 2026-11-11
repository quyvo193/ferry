# M1: Raw TCP Send/Receive + Throughput Baseline

**Budget:** 8 h · **When:** Week 1, sessions 1.1–1.3 · **Tag:** `m1`

## Goal

Send one file from machine A to machine B over plain TCP using a manually entered `ip:port`,
and know how fast the network *can* go (iperf3) vs how fast your code goes.

```sh
B$ ferry listen --port 47701 --dir ./in
A$ ferry send 192.168.1.20:47701 ./big.bin
```

## Concepts to learn

- `net.Listen`, `Listener.Accept`, `net.Dial`, `net.Conn`; the `*net.TCPConn` type
- TCP is a **byte stream**: no message boundaries; `Read` can return fewer bytes than asked
- `io.Copy` / `io.CopyBuffer`, and how they use `io.ReaderFrom` / `io.WriterTo` (e.g. `sendfile`
  on Linux when copying `*os.File` → `*net.TCPConn`). Later framing/TLS disables this fast path.
  Notice that tradeoff.
- `TCPConn.CloseWrite()` (half-close) to signal "no more data" while still being able to read
- `io.ReadFull` for fixed-size reads
- Socket buffers: `SetReadBuffer` / `SetWriteBuffer`; Nagle (`SetNoDelay`) doesn't matter for bulk
- Units: MB vs MiB, MB/s vs Mbit/s (1 Gbit Ethernet ≈ 117 MB/s of payload, best case)

## Tasks

- [ ] Push skeleton to GitHub, confirm CI green on 3 OSes (0.5 h)
- [ ] `flag.FlagSet` for `listen` (`--port`, `--dir`) and `send` (`<ip:port> <file>`) (0.5 h)
- [ ] Receiver: accept loop, one goroutine per conn, create file in `--dir` (1.5 h)
- [ ] Sender: open file, dial, `io.CopyBuffer`, half-close, wait for receiver to close (1 h)
- [ ] Mini header before data so the receiver knows the filename + size
      (e.g. `[u16 name_len][name][u64 size]`). This is a precursor to M2 framing (1 h)
- [ ] Timing on both sides: bytes, duration, MB/s, Mbit/s (0.5 h)
- [ ] `--discard` debug flag on the receiver: write to `io.Discard` to isolate network vs disk (0.25 h)
- [ ] Loopback test in `transfer` package with `t.TempDir()` (1 h)
- [ ] Generate test data: 1 × 4 GiB random file, 10 000 × 4 KiB files (script in `scripts/`, 0.25 h)
- [ ] Buffer sweep 32 KiB / 256 KiB / 1 MiB / 4 MiB on loopback + Ethernet + Wi-Fi (1.5 h)
- [ ] iperf3 baseline between the same two machines; fill BENCHMARKS.md (0.5 h)

## Design decisions

1. **Where does the header live?** Throwaway in M1, real framing in M2. Keep it in `transfer`,
   don't over-design.
2. **Default download dir**: current dir vs `~/Downloads/ferry`. Pick one now and use it consistently.
3. **One conn = one file?** Yes for M1. Multi-file is M2.

## Pitfalls

- Closing the conn on the sender before the receiver has read everything → truncated file.
  Half-close + wait for the receiver's EOF/close fixes this.
- Ignoring the `n` returned by `Read`/`Write`.
- **Page cache lies**: the 2nd run of the same file reads from RAM. Note which run you measure.
- **Disk can be the bottleneck**, not the network. That's what `--discard` is for.
- macOS/Windows firewall prompt on first `listen`. Allow it, or the sender just times out.
- Wi-Fi numbers vary a lot. Take the median of 3 runs and note distance/band.

## Tests

- Loopback: send a random 10 MiB file, compare bytes/hash of received file
- Zero-byte file
- Filename with spaces / unicode

## Done when

- [ ] 4 GiB file A→B over Ethernet between two real machines; `shasum -a 256` matches on both sides
- [ ] BENCHMARKS.md has iperf3 ceiling + buffer sweep; you know which buffer size you'll use
- [ ] You can explain where time goes (network vs disk vs CPU)

## Stretch

- `pprof` CPU profile of a transfer (`runtime/pprof`): see what the hot path actually is
- Compare `io.Copy` (sendfile path) vs manual read/write loop on Linux
