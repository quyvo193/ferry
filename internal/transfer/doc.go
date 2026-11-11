// Package transfer drives a send or receive session on top of package
// protocol: walking files and folders, streaming data with large buffers,
// hashing on both sides, safe writes (file.part -> rename after verify),
// resume from existing offsets and optional parallel streams.
//
// Milestone: M1-M5. Plan: docs/plan/03-m2-framing-protocol.md, 06-m5-resume-parallel-progress.md
package transfer
