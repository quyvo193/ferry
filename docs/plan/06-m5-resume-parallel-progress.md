# M5: Resume, Parallel Streams, Small-File Batching, Progress

**Budget:** 18 h · **When:** Week 3 session 3.5 + Week 4 sessions 4.1–4.4 · **Tag:** `m5`

| Part | Hours | Priority |
|---|---|---|
| A. Resume | 6 | Must |
| B. Progress | 4 | Must |
| C. Small-file batching | 2 | Should |
| D. Parallel streams | 6 | **Could**, time-boxed, cut first |

---

## A. Resume (6 h)

### Goal
Kill a transfer at 60 %, run the same `ferry send` again → it continues from ~60 %, final hash still verified.

### Concepts
- `os.File.Seek`, `os.File.Stat`, `Truncate`
- Why a `.part` file's size isn't proof its content is right (crash before fsync, different source file)

### Tasks
- [ ] Receiver: on OFFER, for each file look for `<path>.part`; decide a resumable offset (2 h)
- [ ] Validate the `.part` really belongs to this file (see decision 1) (1.5 h)
- [ ] ACCEPT carries per-file offsets; sender seeks and starts DATA at that offset (1 h)
- [ ] Hash continuity: receiver re-hashes the existing prefix; sender reads + hashes its prefix
      locally (no send) so both whole-file digests still match (1 h)
- [ ] Fault test: conn wrapper that dies after N bytes → rerun → resumed + verified (0.5 h)

### Design decisions
1. **Is this `.part` the same file?**
   - (a) Sidecar `.part.meta` with sender device_id, path, size, mtime. Simple, cheap.
   - (b) Receiver sends hash of its prefix, sender compares with its own prefix hash. Robust, costs a read on both sides.
   Recommendation: (a). The final full-file hash catches anything that slips through.
2. **Trust the tail?** Round the offset down (e.g. to a 1 MiB boundary) to drop a possibly torn last write.
3. **Resume across different senders** of the same path: no, require matching device_id.

### Pitfalls
- Off-by-one between offset and what was actually fsynced
- Sender's file changed since the first attempt → size/mtime mismatch → restart from 0, don't corrupt

---

## B. Progress (4 h)

### Goal
```
photos/IMG_2041.jpg   [=========>          ]  48.2%  1.21 GB / 2.51 GB  108.4 MB/s  ETA 0:12
```

### Concepts
- `sync/atomic` (`atomic.Int64`) for counters updated by I/O goroutines and read by the renderer
- Counting `io.Writer`/`io.Reader` wrapper
- Exponential moving average for smooth speed; ETA = remaining / speed
- Terminal: `\r` to rewrite the line, clearing leftovers; detecting TTY
  (`os.Stdout.Stat()` → `Mode()&os.ModeCharDevice`)

### Tasks
- [ ] `progress.Tracker`: total, done (atomic), per-file label; `Add(n)` (1 h)
- [ ] Renderer goroutine, 4–5 Hz ticker, EMA speed, ETA, humanized sizes (1.5 h)
- [ ] Non-TTY mode: one log line every few seconds (pipes/CI) (0.5 h)
- [ ] Wire into sender + receiver; final summary line (total, avg speed, duration) (1 h)

### Pitfalls
- Rendering from the I/O goroutine slows the hot path. Keep rendering on its own goroutine.
- Resume: start the bar at the resumed offset, but compute speed only from new bytes.
- Windows terminals: `\r` works; ANSI clear-line works on Windows 10+ terminals.

---

## C. Small-file batching (2 h)

### Goal
10 000 × 4 KiB files transfer at a reasonable fraction of big-file throughput.

### Tasks
- [ ] Benchmark the `small` set first; profile (1 h)
- [ ] Likely fixes: no per-file round trips (there shouldn't be any already), large `bufio.Writer`
      so many small DATA/DONE frames coalesce into big writes, open the next file while sending
      the current one; receiver avoids per-file fsync stalls (decision: fsync per file vs at end) (1 h)

---

## D. Parallel streams (6 h, time-boxed)

### Step 0: decide if it's worth it (30 min)
Compare single-stream TLS throughput vs `iperf3 -P 1` and `-P 4`:
- Single stream ≥ 90 % of the ceiling → **skip D**, write the finding in BENCHMARKS.md. That's a valid result.
- Single stream is CPU-bound (one core at 100 % for TLS/hash) or Wi-Fi gains from `-P 4` → continue.

### Concepts
- `os.File.WriteAt` / `ReadAt` (concurrent positional I/O on one file)
- Preallocate with `Truncate(size)`
- Coordinating N goroutines: first error cancels all (`context.WithCancelCause` + `WaitGroup`)

### Tasks
- [ ] Negotiate: OFFER says "can do N streams", ACCEPT says how many (1 h)
- [ ] Secondary connections: TLS (same pinning) + new `JOIN(transfer_id, stream_index)` message (1.5 h)
- [ ] Split large files (> threshold, e.g. 256 MiB) into ranges; small files stay on stream 0 (1.5 h)
- [ ] Receiver: preallocate `.part`, `WriteAt` per DATA `offset` (0.5 h)
- [ ] Integrity with out-of-order data (see decision) (1 h)
- [ ] `--streams N` flag, default 1 until benchmarks say otherwise (0.5 h)

### Design decisions
1. **Hashing out-of-order data**: streaming whole-file hash no longer works.
   - (a) Re-hash the finished `.part` on the receiver (extra full read, simple)
   - (b) Per-range hashes in DONE, final = hash of range hashes (no extra read, protocol change)
2. **Resume + parallel together**: per-range offsets get complex. It's OK to say "parallel
   transfers resume from 0" in v1.

---

## Done when (M5)

- [ ] Interrupted transfer resumes and verifies (automated test + manual on 2 machines)
- [ ] Progress bar on both sides, sane in non-TTY mode
- [ ] `small` set benchmark recorded before/after batching work
- [ ] Parallel streams either shipped behind `--streams` or consciously cut with numbers to justify it
