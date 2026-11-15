#!/usr/bin/env bash
set -euo pipefail

usage() {
	echo "usage: $0 [big|medium|small|all] [out_dir]" >&2
	echo "  big:    1 x 4 GiB random file" >&2
	echo "  medium: 100 x 40 MiB random files" >&2
	echo "  small:  10000 x 4 KiB random files" >&2
	exit 2
}

set_name=${1:-all}
out=${2:-testdata/large}

rand_file() {
	head -c "$2" /dev/urandom >"$1"
}

gen_big() {
	mkdir -p "$out"
	echo "big: $out/big.bin"
	rand_file "$out/big.bin" $((4 * 1024 * 1024 * 1024))
}

gen_medium() {
	mkdir -p "$out/medium"
	echo "medium: $out/medium/"
	for i in $(seq -w 1 100); do
		rand_file "$out/medium/m$i.bin" $((40 * 1024 * 1024))
	done
}

gen_small() {
	mkdir -p "$out/small"
	echo "small: $out/small/"
	head -c $((10000 * 4096)) /dev/urandom | split -a 5 -b 4096 - "$out/small/s"
}

case "$set_name" in
big) gen_big ;;
medium) gen_medium ;;
small) gen_small ;;
all) gen_big && gen_medium && gen_small ;;
*) usage ;;
esac

du -sh "$out"
