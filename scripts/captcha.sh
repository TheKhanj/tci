#!/bin/bash

dir="captchas"

fetch_captcha() {
	curl --fail --silent --show-error \
		--url 'https://internet.tci.ir/panel/captcha/'
}

fetch() {
	local count=100 help=0

	while getopts "c:h:" opt; do
		case "$opt" in
		c) count="$OPTARG" ;;
		h) help=1 ;;
		*) exit 1 ;;
		esac
	done

	# show help
	if [ "$help" ] ...

	mkdir -p "$dir"

	for ((i = 0; i < count; i++)); do
		id=$(printf "%04d" "$(($(find "$dir" -name '*.png' | wc -l) + 1))")
		local file="$dir/$id.png"

		[ -f "$file" ] && continue

		fetch_captcha >"$file" || {
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

		local vfile="${file%.png}"
		[ -f "$vfile" ] && continue

		local win=$(xdotool getwindowfocus)

		feh "$file" &
		local pid=$!

		xdotool windowfocus --sync "$win"
		read -r -p "Enter the CAPTCHA value: " value

		kill "$pid" 2>/dev/null
		wait "$pid" 2>/dev/null

		[ -n "$value" ] && printf '%s\n' "$value" >"$vfile"
	done
}

check() {
	start=1

	while getopts "s:h:" opt; do
		case "$opt" in
		s) start=$OPTARG ;;
		*) exit 1 ;;
		esac
	done

	# show help
	if [ "$help" ] ...

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
		read -r -p "$id: $value [Y/n/q] " answer

		kill "$pid" 2>/dev/null
		wait "$pid" 2>/dev/null

		case "$answer" in
		n | N)
			rm -f "$vfile"
			echo "Removed label: $id"
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
