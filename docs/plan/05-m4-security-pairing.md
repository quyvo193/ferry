# M4: TLS, Identity, Pairing (TOFU + PIN), Accept Prompt

**Budget:** 16 h · **When:** Week 3, sessions 3.1–3.4 · **Tag:** `m4`

## Goal

All session traffic is encrypted and mutually authenticated. The first connection between two
devices is confirmed with a PIN; after that the peer's certificate is pinned. The receiver must
approve every transfer after seeing who is sending what.

```
$ ferry listen
New device "Quy's MacBook" (a1b2c3d4) wants to pair.
Confirm this PIN matches the sender's screen: 482 913   [y/N] y
Quy's MacBook wants to send 3 items (1.24 GB):
  photos/           (212 files, 1.20 GB)
  notes.txt         4.1 KB
  video.mp4         38.0 MB
Accept? [y/N] (60s) y
```

## Threat model (write this into PROTOCOL.md §10)

- **In scope**: passive eavesdropper on the LAN; active MITM on the LAN (ARP spoofing, fake
  discovery announcements); a stranger sending you files unsolicited.
- **Out of scope (v1)**: compromised endpoint; protecting the private key on disk beyond file perms.

## Concepts to learn

- TLS 1.3 handshake at a high level; why "self-signed" ≠ "insecure" if you **pin** the cert
- `crypto/ecdsa` (P-256) or `crypto/ed25519`; `x509.CreateCertificate`, `encoding/pem`
- Certificate fingerprint = SHA-256 over the DER bytes (`cert.Raw`)
- `tls.Config`: `Certificates`, `MinVersion: tls.VersionTLS13`, `ClientAuth: tls.RequireAnyClientCert`,
  `InsecureSkipVerify` (client), **`VerifyConnection`** (your pinning hook)
- Why `VerifyConnection` over `VerifyPeerCertificate`: the latter isn't called on resumed sessions
- `tls.Server` / `tls.Client` wrapping an existing `net.Conn`; `ConnectionState()`
- TOFU (SSH `known_hosts`), Short Authentication Strings (Bluetooth numeric comparison, Signal safety numbers)

## Tasks

**Identity (3.5 h)**
- [ ] Generate key + self-signed cert on first run (long validity, `NotBefore` a bit in the past);
      write PEM files (key `0600`); load on later runs (2.5 h)
- [ ] Fingerprint (SHA-256 of DER) + human format (e.g. grouped hex) (0.5 h)
- [ ] Device ID derived from fingerprint; update discovery to use it + include fingerprint (0.5 h)

**TLS + pinning (5 h)**
- [ ] Wrap listener/dialer with TLS 1.3, mutual auth; confirm the transfer still works (2 h)
- [ ] `trust` store: `known_peers.json` load/save (atomic write: temp file + rename) (1.5 h)
- [ ] `VerifyConnection` hook: known + matching → OK; known + mismatch → hard fail with an
      SSH-style "IDENTITY CHANGED" error; unknown → mark connection as "needs pairing" (1.5 h)

**Pairing + accept (4.5 h)**
- [ ] PIN = derive from both fingerprints, order-independent (e.g. hash of sorted pair → 6 digits) (1 h)
- [ ] Pairing UX on both sides; store peer only after confirmation (1.5 h)
- [ ] `Approver` interface in `transfer`; CLI implementation prompts on stdin with timeout →
      ACCEPT or REJECT(reason); file-list summary (folders collapsed, sizes humanized) (2 h)

**Verify + docs (3 h)**
- [ ] HELLO check: device_id in HELLO must match the one derived from the TLS cert (0.5 h)
- [ ] Tests: pinned OK, mismatched cert refused, unknown → pairing path, reject path, prompt timeout (1.5 h)
- [ ] Throughput TLS vs plain; check negotiated cipher (AES-GCM vs ChaCha20) → BENCHMARKS.md (0.5 h)
- [ ] PROTOCOL.md §10 security + flow update (0.5 h)

## Design decisions

1. **Key type**: ECDSA P-256 (universally supported) vs Ed25519 (simpler, smaller, fine in Go TLS 1.3).
2. **PIN UX**:
   - *Numeric comparison*: both screens show the PIN, user presses `y` on both. Simple, but users
     tend to press `y` without looking.
   - *Entry*: receiver shows the PIN, sender types it. The sender computes its own PIN and compares.
     Stronger. Recommendation: entry.
3. **Where does pairing happen**: inside `VerifyConnection` (blocks the handshake while the user
   types) vs after the handshake, before HELLO completes (cleaner). Recommendation: handshake accepts
   any cert, record "unknown", then do pairing at the protocol level before OFFER.
4. **Concurrent offers while a prompt is open**: queue vs immediate REJECT("busy"). Recommendation: busy.
5. **Auto-accept from trusted peers**: `listen --auto-accept` flag? Convenient; keep it off by default.
6. **Trust management**: `ferry trust list|remove` now, or edit the JSON by hand for v1?

## Pitfalls

- `InsecureSkipVerify: true` **without** your own verification = encryption without
  authentication = MITM-able. Write a test that proves a different cert is refused.
- Pinning only on one side. Both sender and receiver must verify each other.
- `0600` means nothing on Windows (ACLs). Note it, don't solve it in v1.
- Reading stdin in the connection goroutine while Ctrl-C handling also wants it: one prompt owner.
- Prompt timeout must send REJECT, not just hang the sender.
- Don't print or log private keys. Don't put the key in the repo test fixtures. Generate in `t.TempDir()`.
- Deleting the config dir = new identity = every peer will see "IDENTITY CHANGED". Document it.

## Tests

- Two in-process nodes with separate temp config dirs:
  - first contact → pairing callback invoked, PIN equal on both sides
  - second contact → no pairing
  - regenerate one identity → connection refused with identity-changed error
- Approver stub: accept, reject, timeout
- Manual: Wireshark capture shows only TLS records; no filenames in clear text

## Done when

- [ ] Fresh devices pair once via PIN; subsequent transfers go straight to the accept prompt
- [ ] Changed identity is refused on both sides with a clear message
- [ ] Rejecting an offer or timing out ends the sender cleanly with a message
- [ ] TLS throughput recorded; within ~10 % of plain TCP on Ethernet (AES-GCM w/ hardware support)

## Stretch

- `ferry trust list|remove`
- Show fingerprint in `ferry peers` with a "trusted" column
