# Benchmarks

Methodology: [plan/07-release-benchmarks.md](plan/07-release-benchmarks.md).
Each number is the **median of 3 runs**. Throughput in **MB/s** (10^6 bytes/s), Mbit/s in parentheses.

## Hardware

| Role | Machine | OS | CPU | Disk | NIC |
|---|---|---|---|---|---|
| A | | | | | |
| B | | | | | |

Network:
- Ethernet: switch model, link speed
- Wi-Fi: router model, band (2.4/5/6 GHz), standard (Wi-Fi 5/6/6E/7), distance

## Test sets

| Name | Content |
|---|---|
| big | 1 × 4 GiB random file |
| medium | 100 × 40 MiB |
| small | 10 000 × 4 KiB |

## Ceiling (iperf3)

| Link | iperf3 single stream | iperf3 -P 4 |
|---|---|---|
| Ethernet | | |
| Wi-Fi | | |

## M1: raw TCP, buffer size sweep (big, sink = io.Discard)

Run with `scripts/sweep.sh <file> [host:port]`. Buffer applies to both sender and receiver.
Loopback: Apple M3, 16 GB, macOS 14 (Darwin 23.6), Go 1.27.3, 2026-11-15.

| Buffer | Ethernet | Wi-Fi | Loopback |
|---|---|---|---|
| default (io.Copy) | | | 1344 |
| 32 KiB | | | 1209 |
| 256 KiB | | | 1168 |
| 1 MiB | | | 2277 |
| 4 MiB | | | 2414 |

## Results by milestone

| Milestone | Config | Set | Ethernet | Wi-Fi | % of iperf3 | Notes |
|---|---|---|---|---|---|---|
| M1 | raw TCP | big | | | | |
| M2 | framed + SHA-256 | big | | | | |
| M2 | framed + SHA-256 | small | | | | |
| M4 | TLS 1.3 | big | | | | |
| M5 | TLS, 4 streams | big | | | | |
| M5 | TLS, batched | small | | | | |

## Notes / findings

- M1 loopback: >= 1 MiB buffers almost double throughput over default/32-256 KiB (~1.2-1.3 GB/s -> ~2.3-2.4 GB/s).
  Loopback is CPU/syscall bound, so fewer, larger reads/writes win. Re-check on Ethernet/Wi-Fi,
  where the link (~117 MB/s on 1 GbE) should cap every buffer size.
- Default path on the receiver with `--discard` goes through `io.Discard.ReadFrom` (8 KiB reads),
  which is likely what holds the default row near the small-buffer rows.
