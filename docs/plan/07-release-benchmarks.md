# Release, Install from GitHub, Benchmarks, Final Docs

**Budget:** 6 h · **When:** Week 4, sessions 4.5–4.6 · **Tag:** `v0.1.0`

## A. Release pipeline (3 h)

Already in the skeleton:
- `go.mod` module path `github.com/quyvo193/ferry` → `go install …/cmd/ferry@latest` works once pushed
- `.goreleaser.yaml`: 6 targets (darwin/linux/windows × amd64/arm64), `CGO_ENABLED=0`, version via ldflags
- `.github/workflows/release.yml`: runs GoReleaser on tag push `v*`
- `.github/workflows/ci.yml`: vet + test (race) on 3 OSes + cross-compile job

### Tasks
- [ ] Install GoReleaser locally; `goreleaser check` + `make snapshot`; inspect `dist/` (0.5 h)
- [ ] Verify `ferry version` shows the right version in both paths:
      ldflags (release binary) and build info (`go install …@v0.1.0`) (0.25 h)
- [ ] Write release notes (features, known limitations, firewall/multicast notes) (0.5 h)
- [ ] `git tag -a v0.1.0 -m "v0.1.0" && git push origin v0.1.0`; watch the Action (0.25 h)
- [ ] From a clean machine/dir: `go install github.com/quyvo193/ferry/cmd/ferry@v0.1.0` (0.25 h)
- [ ] Download each OS archive you can test; run `ferry version`; macOS quarantine + Windows
      SmartScreen notes in README (1 h)
- [ ] Smoke test: macOS ↔ Windows or Linux transfer with release binaries (0.25 h)

### Pitfalls
- Tags are immutable in the Go module proxy. **Never re-push a changed `v0.1.0`**; bump to `v0.1.1`.
- `go install` caches through `proxy.golang.org`. A brand new tag may take a minute to appear.
- The module path must match the GitHub repo path exactly (case matters).

## B. Benchmarks (2 h)

### Method
- Same two machines for every run; record hardware/OS/NIC in BENCHMARKS.md
- Data sets: `big` (1 × 4 GiB random), `medium` (100 × 40 MiB), `small` (10 000 × 4 KiB)
- Ceiling: `iperf3 -c <B>` and `iperf3 -c <B> -P 4`
- Each config: 3 runs, report median; note whether the receiver wrote to disk or discarded
- Matrix:

| | Ethernet | Wi-Fi (5 GHz, same room) | Wi-Fi (other room), optional |
|---|---|---|---|
| iperf3 | ✓ | ✓ | |
| ferry, 1 stream, big | ✓ | ✓ | |
| ferry, 1 stream, small | ✓ | ✓ | |
| ferry, N streams, big (if shipped) | ✓ | ✓ | |

### Questions to answer in "Notes / findings"
- What % of iperf3 does ferry reach? What's the bottleneck (CPU/TLS, hashing, disk, network)?
- How much did TLS cost vs M1 raw TCP?
- Does Wi-Fi benefit from parallel streams more than Ethernet? Why might that be?
- Small files: files/s, and where the time goes

## C. Final docs (1 h)

- [ ] `PROTOCOL.md`: no `TBD` left; byte-level examples for HELLO and one DATA frame
- [ ] README: real `ferry peers` / prompt / progress output; roadmap ticked; known limitations
- [ ] Dev log → short retrospective: estimates vs actual per milestone, what you'd do differently

## Done when

- [ ] GitHub release `v0.1.0` with 6 archives + `checksums.txt`
- [ ] `go install github.com/quyvo193/ferry/cmd/ferry@latest` gives a working binary
- [ ] BENCHMARKS.md complete for Ethernet + Wi-Fi with conclusions
