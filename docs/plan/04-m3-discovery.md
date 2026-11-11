# M3: UDP Multicast Discovery

**Budget:** 10 h · **When:** Week 2, sessions 2.3–2.5 · **Tag:** `m3`

## Goal

Devices running `ferry listen` find each other automatically. `ferry peers` lists them and
`ferry send <name> …` works without typing an IP.

```
$ ferry peers
NAME              DEVICE ID   ADDRESS              LAST SEEN
Quy's MacBook     a1b2c3d4    192.168.1.20:47701   1s ago
Office-PC         9f8e7d6c    192.168.1.31:47701   3s ago
```

## Concepts to learn

- UDP: datagrams, no connection, no delivery guarantee, ~1500 B MTU (stay well below)
- Multicast: groups, IGMP join, TTL=1 (stay on the local subnet), admin-scoped `239.0.0.0/8`
- `net.ListenMulticastUDP(network, iface, groupAddr)`, `net.DialUDP` / `WriteToUDP`
- `net.Interfaces()`, `Interface.Flags` (`FlagUp`, `FlagMulticast`, `FlagLoopback`), `Addrs()`
- (Optional dep) `golang.org/x/net/ipv4.PacketConn`: per-interface join, `SetMulticastLoopback`,
  `SetMulticastInterface`. Read why the stdlib API is limited before deciding to use it.
- `time.Ticker`, `context.Context`, `sync.RWMutex`-protected map
- Broadcast (`255.255.255.255` / subnet broadcast) as a fallback when multicast is blocked

## Announcement packet

Keep it small (< 512 B), versioned, and parseable without trusting it:

```
magic "FRRY" | u8 version | u8 kind (ANNOUNCE / QUERY / BYE)
str device_id | str name | u16 tcp_port | str cert_fingerprint (empty until M4)
```

Use the **packet's source IP** as the peer address, not a self-reported IP.

## Tasks

- [ ] Spike (throwaway): send/recv one multicast packet between two real machines; try
      Wi-Fi + Ethernet; note what breaks (1.5 h)
- [ ] Packet encode/decode + tests + fuzz target (1 h)
- [ ] Announcer: send immediately on start, then every 2 s ± jitter; on all suitable interfaces (1.5 h)
- [ ] Listener: parse, drop own device_id, drop malformed/oversized, update table (1 h)
- [ ] Peer table: `map[deviceID]Peer` with `LastSeen`, sweeper removes after TTL; `Snapshot()` (1 h)
- [ ] `BYE` on graceful shutdown (Ctrl-C) → immediate removal (0.5 h)
- [ ] `ferry peers [--wait 3s]`: send `QUERY` so listeners answer right away, collect, print table (1.5 h)
- [ ] `ferry send <peer>`: resolve by exact name → device-ID prefix → literal `ip:port`;
      ambiguous match → error listing candidates (1 h)
- [ ] `config` package: device name (default `os.Hostname()`), `--name` flag (0.5 h)
- [ ] Update PROTOCOL.md §9 (0.5 h)

## Design decisions

1. **Stdlib only vs `x/net/ipv4`**: stdlib `ListenMulticastUDP` joins on one interface (or the
   default). `x/net/ipv4` gives real per-interface control. It's quasi-stdlib (Go team maintained).
   Spike first. Only add it if you hit a real problem.
2. **Device ID**: random once and persisted, or derived from the cert fingerprint (M4)?
   Deriving it (e.g. first 8 bytes of fingerprint, hex) means one less thing to store, and an ID that
   can't be claimed without the key. Recommendation: use a placeholder now, switch to derived in M4.
3. **Does `peers` need a running `listen`?** No. `peers` is a short-lived listener + QUERY. Both
   processes on one machine must be able to bind the discovery port → check address reuse on each OS.
4. **Binary vs JSON packet**: binary is consistent with the main protocol. JSON is easier to eyeball in Wireshark.

## Pitfalls

- **Multiple interfaces** (macOS: `en0`, `utun*` VPN, `bridge*`, `awdl0`; Docker bridges on Linux):
  packets go out the "default" one only. Iterate interfaces, skip down/loopback/point-to-point.
- **Two instances, same port, same host**: address reuse behaves differently on Linux/macOS/Windows.
  Test this early since it's needed for local testing.
- **Router blocks multicast** (AP/client isolation, IGMP snooping, guest networks). Fallback: broadcast,
  and `ferry send ip:port` always works.
- **Firewalls**: Windows Defender blocks inbound UDP by default for unknown apps.
- Your own announcements loop back to you. Filter by device_id.
- Discovery is **unauthenticated**. Anyone can announce any name/fingerprint. Treat it as a hint;
  M4's TLS pinning is the real check.
- Don't log every received packet at info level (2 s × N peers gets noisy).

## Tests

- Packet codec round-trip + fuzz (garbage must never panic)
- Peer table with an injectable clock: add, refresh, expire, BYE
- Resolution: exact name, prefix, ambiguous, ip:port passthrough
- Manual: 2 machines × (Ethernet, Wi-Fi); macOS ↔ Windows at least once

## Done when

- [ ] Two machines see each other within ~3 s of starting `listen`
- [ ] A killed peer disappears within TTL; Ctrl-C'd peer disappears immediately (BYE)
- [ ] `ferry send "<name>" ./folder` works end-to-end
- [ ] Behavior on a multicast-blocked network is documented (and the manual IP path works)

## Stretch

- mDNS/DNS-SD (`_ferry._tcp.local`) for interop with standard tools. Bigger topic, post-v1.
- IPv6 multicast (`ff02::…` link-local)
