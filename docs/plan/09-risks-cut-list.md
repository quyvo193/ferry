# Risks, Cut List, Stretch Goals

## Risks

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Multicast blocked on home/office network (AP isolation, IGMP snooping) | Medium | M3 demo fails | Spike early (session 2.3); broadcast fallback; manual `ip:port` always works |
| Multi-interface / VPN / Docker weirdness on macOS/Windows | High | Peers not found | Iterate interfaces; test with VPN off first; document |
| Only one physical machine available for testing | ? | Can't test LAN/Wi-Fi | Line up a 2nd device (old laptop, Raspberry Pi, a friend's PC) **before week 1** |
| M2 overruns (protocol design rabbit holes) | High | Pushes everything | Decisions section: pick + document + move on; framing must be done by end of session 2.2 |
| TLS/pairing UX takes longer than 16 h | Medium | M5 squeezed | Cut list items 2–3 |
| Windows quirks (paths, firewall, file locking, `Rename` over open file) | Medium | Cross-OS bugs late | CI on Windows from day 1; one manual Windows run per milestone |
| Disk is the bottleneck, not network | Medium | Misleading benchmarks | `--discard` mode from M1; note disk type |
| Learning detours (reading too much) | Medium | Slow progress | 1 h time-box on reading, then write a spike |

## Cut list (cut in this order when behind)

1. **Parallel streams** (M5-D): save ~6 h. Keep the `offset` field in DATA so it can come later.
2. **PIN entry → numeric comparison**: save ~1 h.
3. **`.part.meta` sidecar → size-only resume + final hash**: save ~1 h. Integrity still guaranteed by the hash.
4. **`BYE` / `QUERY` discovery packets**: save ~1 h. `peers` just waits one announce interval.
5. **Small-file batching work**: save ~2 h. Ship with a known-limitation note.
6. **Non-TTY progress mode**: save ~0.5 h.

**Never cut**: path sanitization, hash verification before rename, TLS pinning, accept prompt.
Those are the correctness and security core.

## Stretch goals (after v0.1.0)

- `ferry trust list|remove`; `--auto-accept` for trusted peers
- BLAKE3 hashing; benchmark vs SHA-256
- mDNS / DNS-SD discovery (`_ferry._tcp.local`)
- IPv6 support
- Send text/clipboard snippets
- Compression (zstd) for compressible data, negotiated in OFFER
- Transfer history log
- Homebrew tap / Scoop bucket
- QUIC transport experiment (`quic-go`) vs TCP+TLS
- GUI (Wails / Fyne) on top of the same `internal/` packages
- Interop experiment: read the LocalSend protocol spec and compare design choices
