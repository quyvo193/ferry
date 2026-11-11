# Ferry: 4-Week Development Plan

**Capacity:** ~20 h/week × 4 weeks = **~80 h**
**Suggested window:** Mon 2026-11-16 → Sun 2026-12-13
**Outcome:** `ferry v0.1.0`, a working CLI for macOS/Windows/Linux, installable with
`go install`, plus `PROTOCOL.md` and Ethernet vs Wi-Fi benchmark notes.

## Files

| File | Content |
|---|---|
| [00-timeline.md](00-timeline.md) | Week-by-week, session-by-session schedule |
| [01-architecture.md](01-architecture.md) | Package layout, data flow, concurrency, on-disk layout, conventions |
| [02-m1-raw-tcp.md](02-m1-raw-tcp.md) | M1: raw TCP send/receive, throughput baseline |
| [03-m2-framing-protocol.md](03-m2-framing-protocol.md) | M2: frame codec, messages, multi-file and folders, hashing, safe writes |
| [04-m3-discovery.md](04-m3-discovery.md) | M3: UDP multicast discovery, peer table |
| [05-m4-security-pairing.md](05-m4-security-pairing.md) | M4: TLS, identity, TOFU pairing, PIN, accept prompt |
| [06-m5-resume-parallel-progress.md](06-m5-resume-parallel-progress.md) | M5: resume, parallel streams, small-file batching, progress |
| [07-release-benchmarks.md](07-release-benchmarks.md) | Release pipeline, `go install`, benchmarks, final docs |
| [08-testing-strategy.md](08-testing-strategy.md) | Unit, fuzz, integration, fault injection, manual matrix |
| [09-risks-cut-list.md](09-risks-cut-list.md) | Risks, what to cut when behind, stretch goals after v0.1 |

## Budget overview

| Block | Hours | Week |
|---|---|---|
| M1: raw TCP + baseline | 8 | 1 |
| M2: framing, multi-file, hashing | 18 | 1–2 |
| M3: discovery | 10 | 2 |
| Buffer / refactor | 4 | 2 |
| M4: TLS + pairing + accept | 16 | 3 |
| M5: resume, parallel, progress | 18 | 3–4 |
| Release + benchmarks + docs | 6 | 4 |
| **Total** | **80** | |

Parallel streams (≈ 6 h in M5) also act as the **second buffer**: if you're behind by week 4, cut it
(see [09-risks-cut-list.md](09-risks-cut-list.md)).

## How to use this plan

- Each milestone file has the same sections: **Goal → Concepts to learn → Tasks (with hour
  estimates) → Design decisions → Pitfalls → Tests → Done when → Stretch**.
- **Design decisions** are open questions on purpose. Pick one, write down *why* (a line in
  `PROTOCOL.md` or a commit message), and move on. Don't stay stuck on them.
- Tick checkboxes as you go. Record actual hours next to the estimate (e.g. `(2h → 3.5h)`) so
  you learn to estimate better.

## Working rules

1. **Ship each milestone end-to-end before polishing.** A milestone is done only when the demo in
   *Done when* works between two real machines, not just on loopback.
2. **Measure before optimizing.** Every perf change gets a before/after number in `BENCHMARKS.md`.
3. **Tests come with the code.** At minimum: codec round-trips, path sanitization, one loopback
   integration test per milestone.
4. **Keep `main` green.** `make check` passes before every push. CI runs on 3 OSes.
5. **Tag milestones**: `m1`, `m2`, … (lightweight tags), then `v0.1.0` at the end.
6. **Spec follows code**: update `docs/PROTOCOL.md` in the same commit that changes the wire format.
7. **Time-box research**: if a concept takes > 1 h of reading without code, write a spike instead.

## Definition of done (v0.1.0)

- [ ] `ferry listen` / `ferry peers` / `ferry send` work between two different OSes on the same LAN
- [ ] All traffic is TLS; first contact requires PIN confirmation; changed identity is refused
- [ ] Receiver approves every transfer after seeing sender name + file list + sizes
- [ ] Folders keep their structure; path traversal is impossible
- [ ] Every file is hash-verified; nothing appears under its final name unless verified
- [ ] An interrupted transfer resumes instead of restarting
- [ ] Progress shows speed, %, ETA
- [ ] Single-stream throughput on Ethernet ≥ 80 % of iperf3
- [ ] `go install github.com/quyvo193/ferry/cmd/ferry@v0.1.0` works; release has 6 binaries
- [ ] `PROTOCOL.md` complete; `BENCHMARKS.md` filled for Ethernet and Wi-Fi
