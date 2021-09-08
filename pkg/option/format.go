package option

import (
	"fmt"
	"strings"
)

func FormatOptionDescription(prefix string, width int, options []Description) string {
	var s strings.Builder

	var w = 0
	for _, opt := range options {
		if ww := len(opt.Option); ww > w {
			w = ww
		}
	}
	var padding = strings.Repeat(" ", w+5)
	var rw = width - len(padding)

	for _, opt := range options {
		if opt.Description == "" {
			fmt.Fprintf(&s, "%v%*v\n", prefix, -w, opt.Option)
		} else {
			fmt.Fprintf(&s, "%v%*v : ", prefix, -w, opt.Option)
			for i, l := range lineWrap(opt.Description, rw) {
				if i != 0 {
					fmt.Fprint(&s, padding)
				}
				fmt.Fprintf(&s, "%v\n", l)
			}
		}
	}
	return s.String()
}

func lineWrap(str string, w int) (lines []string) {
	b, ws, we := nextWordRange(str, 0)
	ls, le := ws, we
	for ws != we {
		if we-ls > w || b {
			lines = append(lines, str[ls:le])
			ls, le = ws, we
		}
		le = we
		b, ws, we = nextWordRange(str, we)
	}

	if le-ls > 0 {
		lines = append(lines, str[ls:le])
	}
	return
}

var asciiSpace = [256]bool{
	'\t': true, '\n': true, '\v': true, '\f': true, '\r': true, ' ': true,
}

func nextWordRange(s string, ii int) (b bool, i, j int) {
	b = false
	i = ii
	l := len(s)
	for i < l && asciiSpace[s[i]] {
		if s[i] == '\r' || s[i] == '\n' {
			b = true
		}
		i++
	}
	j = i
	for j < l && !asciiSpace[s[j]] {
		j++
	}
	return
}

func FormatCompletion(width int, options []Description) string {
	var s strings.Builder

	var w = 0
	for _, opt := range options {
		if ww := len(opt.Option); ww > w {
			w = ww
		}
	}
	var dw = width - w - 3
	if dw < 3 {
		dw = 0
	}

	for _, opt := range options {
		if dw == 0 || opt.Description == "" {
			fmt.Fprintf(&s, "%*v\n", -w, opt.Option)
		} else {
			var d = opt.Description
			var dl = len(opt.Description)
			if dl > dw {
				d = d[:dw-3] + "..."
			}
			fmt.Fprintf(&s, "%*v : %v\n", -w, opt.Option, d)
		}
	}
	return s.String()
}
