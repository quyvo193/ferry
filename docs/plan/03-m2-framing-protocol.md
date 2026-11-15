# M2: Framing Protocol, Multi-File, Hashing, Safe Writes

**Budget:** 18 h · **When:** Week 1 sessions 1.4–1.6 (12 h) + Week 2 sessions 2.1–2.2 (6 h) · **Tag:** `m2`

## Goal

Replace the M1 ad-hoc header with a real, versioned, documented binary protocol.
Send multiple files and folders (structure preserved), hash both sides, and never leave a
corrupt file under its final name.

## Concepts to learn

- `encoding/binary`: `BigEndian.Uint32`, `BigEndian.AppendUint32`, `binary.Read/Write`
  (convenient but slower; note why)
- `bufio.Reader` / `bufio.Writer`: buffering, `Flush`, why mixing buffered + raw reads breaks
- `io.ReadFull`, `io.LimitReader`, `io.TeeReader`, `io.MultiWriter`
- `hash.Hash`, `crypto/sha256`: streaming hashing
- `path/filepath`: `WalkDir`, `Rel`, `ToSlash`/`FromSlash`, `IsLocal`
- `os.Root` (Go 1.24+): file operations confined to a directory, the cleanest anti-traversal tool
- `os.File.Sync` (fsync) and `os.Rename` atomicity (same filesystem only)
- Protocol state machines: which message is legal in which state

## Reading

⭐ = read before coding. The rest you can look up when you reach the task.

### Framing & binary encoding

- ⭐ [Length-prefix framing](https://eli.thegreenplace.net/2011/08/02/length-prefix-framing-for-protocol-buffers) (Eli Bendersky): why TCP needs framing, and the exact `[len][payload]` pattern
- ⭐ [RFC 9113 §4.1 Frame Format](https://www.rfc-editor.org/rfc/rfc9113#section-4.1) (HTTP/2): a real `[length][type][flags][payload]` frame. Also look at §5.1 *Stream States* for how a protocol describes which frame is legal when
- [Beej's Guide §7.5–7.6](https://beej.us/guide/bgnet/html/#serialization): *Serialization* and *Data Encapsulation*, i.e. how to pack fields into bytes and frame them with a length header
- [The byte order fallacy](https://commandcenter.blogspot.com/2012/04/byte-order-fallacy.html) (Rob Pike): why `binary.BigEndian.Uint32(b)` is enough and you never need to know the host's endianness
- [`encoding/binary`](https://pkg.go.dev/encoding/binary): read `ByteOrder`, `AppendByteOrder`. `binary.Read/Write` use reflection, which is why they're slower
- [`bufio`](https://pkg.go.dev/bufio): `Reader`, `Writer`, `Flush`. Once a `bufio.Reader` wraps the conn, read **only** through it, because it may already hold bytes you want
- [`io`](https://pkg.go.dev/io): `ReadFull`, `LimitReader`, `TeeReader`, `MultiWriter`

### Real protocols to steal ideas from (skim, don't copy)

- ⭐ [Syncthing BEP v1](https://docs.syncthing.net/specs/bep-v1.html): a LAN file-sync protocol with hello, length-prefixed messages, file index, block hashes. Closest to ferry
- [LocalSend protocol](https://github.com/localsend/protocol): a LAN file-transfer app (HTTP-based). Useful for the OFFER/ACCEPT flow and how they describe files

### Hashing

- [`hash`](https://pkg.go.dev/hash) + [`crypto/sha256`](https://pkg.go.dev/crypto/sha256): a `hash.Hash` is an `io.Writer`, so `io.MultiWriter(file, hasher)` hashes as you write

### Paths & traversal

- ⭐ [Traversal-resistant file APIs](https://go.dev/blog/osroot) (Go blog): `filepath.IsLocal`, `os.Root`, and why string checks alone aren't enough
- [Zip Slip](https://security.snyk.io/research/zip-slip-vulnerability): the same bug class in archive extractors. It's exactly the attack a malicious OFFER would try
- [Windows: Naming files](https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file): reserved names (`CON`, `NUL`, `COM1`…), invalid characters, `MAX_PATH`
- [`path/filepath`](https://pkg.go.dev/path/filepath): `WalkDir`, `Rel`, `ToSlash`, `FromSlash`

### Safe writes (fsync + rename)

- ⭐ [Files are hard](https://danluu.com/file-consistency/) (Dan Luu): why "write then rename" isn't automatically safe
- [Ensuring data reaches disk](https://lwn.net/Articles/457667/) (LWN): what `fsync` really does. Includes the "also fsync the directory" detail
- [google/renameio](https://github.com/google/renameio): a small Go implementation of temp file → fsync → rename. Short enough to read the whole source

### Fuzzing

- [Go fuzzing tutorial](https://go.dev/doc/tutorial/fuzz) then the [reference](https://go.dev/doc/security/fuzz/)

## Wire format (starting point, record final choices in PROTOCOL.md)

```
Frame:  [u32 length][u8 type][payload]
```

| Message | Payload idea |
|---|---|
| HELLO | magic `"FRRY"`, `u16 version`, `str device_id`, `str name` |
| OFFER | `u32 count`, then per entry: `str path`, `u8 kind (file/dir)`, `u64 size`, `u32 mode` |
| ACCEPT | `u32 count`, per file `u64 start_offset` (all 0 until M5) |
| REJECT | `u16 reason`, `str message` |
| DATA | `u32 file_index`, `u64 offset`, raw bytes |
| DONE | `u32 file_index`, `u8 hash_alg`, `[32]byte digest`; special index for "transfer complete" |
| ERROR | `u16 code`, `str message` |

## Tasks

**Protocol package (week 1)**
- [ ] Frame writer: write header + payload, reusing buffers (no alloc per DATA frame) (1 h)
- [ ] Frame reader: read header, enforce max size, read payload into a reusable buffer (1 h)
- [ ] Message structs + encode/decode for each type (2 h)
- [ ] Round-trip table tests for every message; truncated-input tests (1 h)
- [ ] HELLO exchange + version check → ERROR on mismatch (1 h)

**Transfer (week 1)**
- [ ] Sender: expand CLI args into an entry list (files + folders via `WalkDir`, empty dirs too,
      skip symlinks with a warning) (2 h)
- [ ] OFFER → ACCEPT/REJECT flow (receiver auto-accepts for now; prompt comes in M4) (1 h)
- [ ] DATA streaming with 1 MiB chunks, `file_index` + `offset`; hashing on both sides while
      streaming; DONE carries sender's digest; receiver compares (3 h)

**Transfer (week 2)**
- [ ] Receiver path safety: reject absolute, `..`, empty, NUL; use `os.Root`; Windows reserved
      names (`CON`, `NUL`, `COM1`…) and invalid chars (`:*?"<>|`) (1.5 h)
- [ ] Safe writes: write `name.part` → `Sync` → verify hash → `Rename`. On mismatch: delete
      `.part` and report ERROR (1 h)
- [ ] Name conflict policy when final file exists (0.5 h)
- [ ] Integration test: folder tree with nested dirs, empty dir, empty file, unicode names (1.5 h)
- [ ] Fuzz test for frame reader + OFFER decoder (`go test -fuzz`) (1 h)
- [ ] PROTOCOL.md v0.1: framing, all messages, flow, errors (1 h)
- [ ] Re-measure throughput vs M1 (framing + hashing cost) → BENCHMARKS.md (0.5 h)

## Design decisions

1. **`length` covers what?** `type + payload` (common) or `payload` only. Either is fine. Document it.
2. **Control message encoding**: hand-rolled binary (learning goal, compact) vs JSON payload
   (easy debug/evolve). DATA must stay binary regardless. Recommendation: binary everywhere,
   since that's the point of the project.
3. **DATA carries `file_index + offset`** or is implicitly "next bytes of current file"?
   Explicit costs 12 bytes per MiB (~0.001 %) and makes resume + parallel streams (M5) much easier.
   Recommended.
4. **Hash scope**: whole-file hash (simple) vs per-chunk hashes (enables parallel + finer resume).
   Recommendation: whole-file now. Revisit in M5 only if you do parallel streams.
5. **Conflict policy**: overwrite / skip / auto-rename `name (1).ext` / reject the whole offer.
   Recommendation: auto-rename. It's the least surprising for users.
6. **Preserve mode/mtime?** Mode bits: yes on Unix, ignore on Windows. mtime: optional stretch.
7. **Unknown message type**: error out (strict) vs skip (forward-compatible). Strict is easier for v1.

## Pitfalls

- Forgetting `bufio.Writer.Flush()` → the peer hangs waiting.
- Allocating `make([]byte, length)` straight from the wire without a max → memory DoS.
- `length == 0` or `length < 1` (no type byte) → must be an error, not a panic.
- Reading the file twice (once to hash, once to send). Use `TeeReader`/`MultiWriter`, one pass.
- File changes size while sending → detect (bytes sent ≠ offered size) and fail clearly.
- Files > 4 GiB: any `uint32` size/offset is a bug.
- macOS filenames may be NFD-normalized unicode. Know it exists, don't fix it in v1.
- Windows `MAX_PATH` (260) for deep trees: note it, test it once.
- `os.Rename` across filesystems fails. `.part` must live next to the final file.

## Tests

- Codec: round-trip every message; decode truncated/oversized/garbage → error, never panic
- Fuzz: `FuzzReadFrame`, `FuzzDecodeOffer`
- Path sanitizer: table test with `../x`, `/etc/passwd`, `C:\x`, `a/../../b`, `CON.txt`, `a\x00b`
- Integration: folder tree A→B on loopback, compare tree + hashes
- Corruption: flip one byte in a DATA payload (test conn wrapper) → hash mismatch → no final file

## Done when

- [ ] `ferry send ip:port ./dir file1 file2` recreates the same tree on B, all hashes verified
- [ ] A malicious OFFER with `../` paths is rejected by an automated test
- [ ] Killing the sender mid-file leaves only `.part`, never a truncated final file
- [ ] PROTOCOL.md describes everything on the wire
- [ ] Throughput within ~10 % of M1 (else profile: allocations? small writes? hashing?)

## Stretch

- BLAKE3 behind a `hash_alg` byte; benchmark vs SHA-256 on your CPUs
- Preserve mtime
