#!/bin/bash

captchasDir="captchas"
labelsDir="labels"
checksDir="checks"

fetch_captcha() {
	curl --fail --silent --show-error \
		--url 'https://internet.tci.ir/panel/captcha/'
}

fetch() {
	local count=100 help=0 opt

	while getopts "c:h" opt; do
		case "$opt" in
		c) count="$OPTARG" ;;
		h) help=1 ;;
		*) return 1 ;;
		esac
	done

	if ((help)); then
		cat <<EOF
Usage: captcha.sh fetch [-c COUNT] [-h]

Fetch new CAPTCHA images.

Options:
  -c COUNT    Number of images to fetch (default: 100)
  -h          Show this help message
EOF
		return 0
	fi

	mkdir -p "$captchasDir" 2>/dev/null

	for ((i = 0; i < count; i++)); do
		id=$(printf "%04d" "$(($(find "$captchasDir" -name '*.png' | wc -l) + 1))")
		local file="$captchasDir/$id.png"

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
	if [[ "$1" == "-h" ]]; then
		cat <<EOF
Usage: captcha.sh label [-h]

Interactively label unlabelled CAPTCHA images.

Controls:
  Enter       Save the entered CAPTCHA value
  Empty       Skip the image
  Ctrl+C      Exit

Options:
  -h          Show this help message
EOF
		return 0
	fi

	mkdir -p "$captchasDir"

	for captcha in "$captchasDir"/*.png; do
		[ -f "$captcha" ] || continue

		local label="${captcha%.png}"
		label="${label/captchas/labels}"
		[ -f "$label" ] && continue

		local win
		win=$(xdotool getwindowfocus)

		feh "$captcha" &
		local pid=$!

		xdotool windowfocus --sync "$win"
		read -r -p "Enter the CAPTCHA value: " value

		kill "$pid" 2>/dev/null
		wait "$pid" 2>/dev/null

		[ -n "$value" ] && printf '%s\n' "$value" >"$label"
	done
}

check() {
	local start=1 help=0 opt

	while getopts "s:h" opt; do
		case "$opt" in
		s) start="$OPTARG" ;;
		h) help=1 ;;
		*) return 1 ;;
		esac
	done

	if ((help)); then
		cat <<EOF
Usage: captcha.sh check [-s ID] [-h]

Interactively verify existing CAPTCHA labels.

Options:
  -s ID       Start checking from this image ID (default: 0001)
  -h          Show this help message

Controls:
  Enter / y   Label is correct
  n           Delete the label
  q           Quit
EOF
		return 0
	fi

	mkdir -p "$labelsDir" 2>/dev/null

	for captcha in "$captchasDir"/*.png; do
		[ -f "$captcha" ] || continue

		local id="${captcha##*/}"
		id="${id%.png}"

		((10#$id < 10#$start)) && continue

		local label="$labelsDir/$id"
		[ -f "$label" ] || continue

		local win
		win=$(xdotool getwindowfocus)

		feh "$captcha" &
		local pid=$!

		xdotool windowfocus --sync "$win"
		local check checkFile="$checksDir/$id"
		check="$(cat "$checkFile" || echo 0)"

		local value answer
		value=$(cat "$label")
		read -r -p "$id: $value (checked $check times) [Y/n/q] " answer

		kill "$pid" 2>/dev/null
		wait "$pid" 2>/dev/null

		case "$answer" in
		n | N)
			rm -f "$label" "$checkFile"
			echo "Removed label: $id"
			;;
		q | Q)
			break
			;;
		esac
		echo "$((check + 1))" >"$checkFile"
	done
}

case "$1" in
fetch)
	shift
	fetch "$@"
	;;
label)
	shift
	label "$@"
	;;
check)
	shift
	check "$@"
	;;
-h)
	cat <<EOF
Usage: captcha.sh <command> [options]

Commands:
  fetch       Fetch CAPTCHA images
  label       Label unlabelled images
  check       Verify existing labels

Run 'captcha.sh <command> -h' for command-specific help.
EOF
	;;
*)
	echo "Usage: captcha.sh <fetch|label|check> [options]" >&2
	exit 1
	;;
esac
