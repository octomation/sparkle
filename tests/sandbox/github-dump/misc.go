package main

import (
	"strings"
	"unicode/utf8"
)

func emoji(in string) (string, string) {
	r, _ := utf8.DecodeRuneInString(in)
	e := string(r)

	// naive implementation
	// research:
	// - https://github.com/spatie/emoji
	// - https://github.com/enescakir/emoji
	emojiRanges := [...]rune{
		0x1F600, 0x1F64F, // Emoticons
		0x1F300, 0x1F5FF, // Misc Symbols and Pictographs
		0x1F680, 0x1F6FF, // Transport and Map
		0x2600, 0x26FF, // Misc symbols
		0x2700, 0x27BF, // Dingbat symbols

		0x1FAA0, 0x1FAA8, // Manually maintained
	}

	for i := 0; i < len(emojiRanges); i += 2 {
		if r >= emojiRanges[i] && r <= emojiRanges[i+1] {
			return e, strings.TrimSpace(strings.TrimPrefix(in, e))
		}
	}
	return "", in
}
