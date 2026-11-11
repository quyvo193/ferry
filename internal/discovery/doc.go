// Package discovery announces this device on the LAN via UDP multicast and
// maintains a table of peers that expire after a TTL.
//
// Announcement: device_id, name, TCP port, certificate fingerprint.
// Discovery data is unauthenticated: it is a hint, never a source of trust.
//
// Milestone: M3. Plan: docs/plan/04-m3-discovery.md
package discovery
