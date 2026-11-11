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

| Buffer | Ethernet | Wi-Fi | Loopback |
|---|---|---|---|
| 32 KiB | | | |
| 256 KiB | | | |
| 1 MiB | | | |
| 4 MiB | | | |

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

- 
