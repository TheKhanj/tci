#!/bin/bash

dir="captchas"

get_captcha() {
	curl --fail --silent --show-error \
		--url 'https://internet.tci.ir/panel/captcha/'
}

fetch() {
	count=100

	while getopts "c:" opt; do
		case "$opt" in
		c) count="$OPTARG" ;;
		*) exit 1 ;;
		esac
	done

	mkdir -p "$dir"

	for ((i = 0; i < count; i++)); do
		id=$(printf "%04d" "$(($(find "$dir" -name '*.png' | wc -l) + 1))")
		file="$dir/$id.png"

		[ -f "$file" ] && continue

		get_captcha >"$file" || {
			rm -f "$file"
			echo "Failed to fetch $id" >&2
			return 1
		}

		echo "Fetched $id"
	done
}

label() {
	mkdir -p "$dir"

	for file in "$dir"/*.png; do
		[ -f "$file" ] || continue

		vfile="${file%.png}"
		[ -f "$vfile" ] && continue

		win=$(xdotool getwindowfocus)

		feh "$file" &
		pid=$!

		xdotool windowfocus --sync "$win"
		read -r -p "Enter the number: " value

		kill "$pid" 2>/dev/null
		wait "$pid" 2>/dev/null

		[ -n "$value" ] && printf '%s\n' "$value" >"$vfile"
	done
}

check() {
	start=1

	while getopts "s:" opt; do
		case "$opt" in
		s) start=$OPTARG ;;
		*) exit 1 ;;
		esac
	done

	for file in "$dir"/*.png; do
		[ -f "$file" ] || continue

		id="${file##*/}"
		id="${id%.png}"

		((10#$id < 10#$start)) && continue

		vfile="$dir/$id"
		[ -f "$vfile" ] || continue

		win=$(xdotool getwindowfocus)

		feh "$file" &
		pid=$!

		xdotool windowfocus --sync "$win"

		value=$(cat "$vfile")
		read -r -p "$id: $value [Y/n/s/q] " answer

		kill "$pid" 2>/dev/null
		wait "$pid" 2>/dev/null

		case "$answer" in
		n | N)
			rm -f "$vfile"
			echo "Removed label: $id"
			;;
		s | S)
			continue
			;;
		q | Q)
			break
			;;
		esac
	done
}

case "$1" in
fetch)
	shift
	fetch "$@"
	;;
label)
	label
	;;
check)
	shift
	check "$@"
	;;
*)
	echo "Usage: captcha.sh {fetch|label|check}"
	exit 1
	;;
esac
