// Package identity owns this device's long-lived identity: a self-signed TLS
// certificate + private key generated once on first run, the derived
// certificate fingerprint, and the device ID.
//
// Milestone: M4. Plan: docs/plan/05-m4-security-pairing.md
package identity
