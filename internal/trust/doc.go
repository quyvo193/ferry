// Package trust implements trust-on-first-use (TOFU) pairing: the store of
// known peers (device ID -> certificate fingerprint), the short PIN used to
// confirm a first pairing, and the TLS verification hook that pins peers.
//
// Milestone: M4. Plan: docs/plan/05-m4-security-pairing.md
package trust
