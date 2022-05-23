package strcase

import (
	"bufio"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Split returns a list of individual words from a string that is either space
// separated or any supported strcase format.
func Split(s string) (r []string) {
	var scanner = bufio.NewScanner(strings.NewReader(s))
	scanner.Split(scanFragment)

	for scanner.Scan() {
		r = append(r, scanner.Text())
	}
	return
}

func scanFragment(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSeparator(r) {
			break
		}
	}

	// Assume we are starting in an uppercase sequence
	var upperStart = start
	var upperEnd = start
	var upperSequence = true

	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])

		if upperSequence {
			if unicode.IsUpper(r) {
				upperEnd = i
			} else {
				upperSequence = false
				if upperStart != upperEnd {
					// Non-empty uppercase sequence; return as a word excluding
					// the last uppercase letter.
					return upperEnd, data[upperStart:upperEnd], nil
				}
			}
		} else if unicode.IsUpper(r) {
			// Not in uppercase sequence, uppercase starts a new fragment.
			return i, data[start:i], nil
		}

		if isSeparator(r) {
			return i + width, data[start:i], nil
		}
	}

	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

func isSeparator(c rune) bool {
	switch c {
	case ' ', '\t', '\n', '\v', '\f', '\r':
		return true
	case '-', '_':
		return true
	}
	return false
}
