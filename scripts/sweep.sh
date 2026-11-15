#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
	echo "usage: $0 <file> [host:port]" >&2
	echo "  no host:port: loopback, starts a local 'ferry listen --discard' per buffer size" >&2
	echo "  host:port:    remote receiver, start it with 'ferry listen --discard [--buf N]'" >&2
	echo "env: BUFS (KiB list, default '32 256 1024 4096'), RUNS (default 3), FERRY (default ./bin/ferry)" >&2
	exit 2
fi

file=$1
addr=${2:-}
ferry=${FERRY:-./bin/ferry}
runs=${RUNS:-3}
bufs=${BUFS:-"32 256 1024 4096"}
port=47799

median() {
	sort -n | awk '{ v[NR] = $1 } END { print (NR % 2) ? v[(NR + 1) / 2] : (v[NR / 2] + v[NR / 2 + 1]) / 2 }'
}

field() {
	sed -nE "s|.*[ ]$1=([0-9.e+]+).*|\1|p"
}

listener=""
cleanup() {
	if [ -n "$listener" ]; then kill "$listener" 2>/dev/null || true; fi
}
trap cleanup EXIT

printf "%-10s %12s %12s   %s\n" "buffer" "MB/s" "Mbit/s" "runs (MB/s)"
for kib in default $bufs; do
	if [ "$kib" = default ]; then bytes=0; else bytes=$((kib * 1024)); fi

	target=$addr
	if [ -z "$addr" ]; then
		"$ferry" listen --port "$port" --discard --buf "$bytes" >/dev/null 2>&1 &
		listener=$!
		sleep 0.3
		target=127.0.0.1:$port
	fi

	mbs=()
	mbits=()
	for _ in $(seq "$runs"); do
		out=$("$ferry" send --buf "$bytes" "$target" "$file" 2>&1)
		mbs+=("$(echo "$out" | field 'MB/s')")
		mbits+=("$(echo "$out" | field 'Mbit/s')")
	done

	if [ -n "$listener" ]; then
		kill "$listener" 2>/dev/null || true
		wait "$listener" 2>/dev/null || true
		listener=""
	fi

	label=$kib
	[ "$kib" != default ] && label="${kib} KiB"
	printf "%-10s %12.1f %12.1f   %s\n" "$label" \
		"$(printf '%s\n' "${mbs[@]}" | median)" \
		"$(printf '%s\n' "${mbits[@]}" | median)" \
		"$(printf '%.0f ' "${mbs[@]}")"
done
