# Ferry Wire Protocol

> **Status: draft.** Sections marked `TBD` get filled in as milestones land
> (M2: framing + messages, M3: discovery packet, M4: TLS + pairing, M5: resume + parallel).

Protocol version: `1` (draft)

## 1. Overview

Ferry has two channels:

| Channel | Transport | Purpose |
|---|---|---|
| Discovery | UDP multicast | Announce presence on the LAN |
| Session | TCP, wrapped in TLS 1.3 (from M4) | Handshake, offer/accept, file data |

All multi-byte integers are **big-endian**. Strings are UTF-8.

## 2. Conventions

- `u8`, `u16`, `u32`, `u64`: unsigned integers, big-endian
- `str`: `[u16 length][bytes]` (TBD: confirm encoding)
- Paths on the wire are relative and use `/` as the separator

## 3. Framing

```
+----------------+------------+---------------------+
| length (u32)   | type (u8)  | payload             |
+----------------+------------+---------------------+
```

- `length`: TBD. Decide whether it covers `type + payload` or only `payload`.
- Maximum frame size: TBD (e.g. 16 MiB). Receivers MUST reject larger frames.
- Unknown `type`: TBD (error vs. ignore)

## 4. Message types

| Type | Value | Direction | Description |
|---|---|---|---|
| HELLO  | TBD | both | Version + identity exchange |
| OFFER  | TBD | sender → receiver | List of files with sizes |
| ACCEPT | TBD | receiver → sender | Approves the offer (+ resume offsets) |
| REJECT | TBD | receiver → sender | Declines the offer, with reason |
| DATA   | TBD | sender → receiver | File bytes |
| DONE   | TBD | sender → receiver | End of file (+ hash) / end of transfer |
| ERROR  | TBD | both | Fatal error, connection closes after |

### 4.1 HELLO
TBD: magic, protocol version, device ID, device name.

### 4.2 OFFER
TBD: transfer ID, entries (path, type file/dir, size u64, mode?, mtime?), totals.

### 4.3 ACCEPT / REJECT
TBD: per-file start offsets (resume), stream count (parallel). REJECT reason codes.

### 4.4 DATA
TBD: file index, offset, bytes.

### 4.5 DONE
TBD: file index + hash algorithm + digest; final DONE for the whole transfer.

### 4.6 ERROR
TBD: error code (u16) + message.

## 5. Session flow

```
Sender                                Receiver
  |  ---- TCP connect + TLS ---------->  |
  |  ---- HELLO ---------------------->  |
  |  <--- HELLO -----------------------  |
  |  ---- OFFER ---------------------->  |
  |                                      |  (user approves)
  |  <--- ACCEPT / REJECT -------------  |
  |  ---- DATA ... DATA -------------->  |
  |  ---- DONE (file 0, hash) -------->  |
  |  ---- DATA ... DONE (file n) ----->  |
  |  ---- DONE (transfer) ------------>  |
  |  <--- result / ERROR --------------  |
```

TBD: final confirmation from the receiver (all hashes verified?).

## 6. Integrity

TBD: hash algorithm (SHA-256 / BLAKE3), what is hashed (whole file vs. chunks), behavior on mismatch.

## 7. Resume

TBD: how the receiver reports offsets, how stale `.part` files are detected.

## 8. Parallel streams

TBD: how secondary connections join a transfer, range assignment.

## 9. Discovery

- Multicast group / port: TBD
- Interval / TTL: TBD

Announcement packet: TBD (magic, version, device ID, name, TCP port, cert fingerprint).

## 10. Security

TBD: TLS settings, certificate pinning, PIN derivation, threat model.
Discovery packets are unauthenticated and are only hints. Trust comes from the TLS certificate fingerprint.

## 11. Versioning

TBD: what happens on version mismatch in HELLO.

## 12. Error codes

| Code | Name | Meaning |
|---|---|---|
| TBD | | |
