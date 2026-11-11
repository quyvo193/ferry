// Package protocol implements ferry's wire format: length-prefixed binary
// frames carried over a TCP (later TLS) connection.
//
// Frame layout, big-endian:
//
//	[u32 length][u8 type][payload]
//
// Messages: HELLO, OFFER, ACCEPT, REJECT, DATA, DONE, ERROR.
// docs/PROTOCOL.md is the authoritative spec.
//
// Milestone: M2. Plan: docs/plan/03-m2-framing-protocol.md
package protocol
