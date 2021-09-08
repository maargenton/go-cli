package option

import "strings"

// scanTagFields scans comma-separated fields in a StructTag, with optional
// value after a ':' separator. Commas and columns can be escaped with '\'
func scanTagFields(str string) (field, value, rest string) {
	var i = 0
	for ; i < len(str); i++ {
		if str[i] == ',' {
			break
		}
		if str[i] == '\\' {
			if i+1 >= len(str) {
				break
			}
			i++
		}
	}

	if i < len(str) {
		rest = str[i+1:]
	} else {
		rest = ""
	}

	field, value = parseTagField(str[:i])
	return
}

func parseTagField(str string) (field, value string) {
	var i = 0
	for ; i < len(str); i++ {
		if str[i] == ':' {
			break
		}
	}

	if i+1 < len(str) {
		return unescapeField(str[:i]), unescapeField(str[i+1:])
	}
	return unescapeField(str), ""
}

func unescapeField(str string) string {
	i := strings.Index(str, "\\")
	if i < 0 {
		return strings.TrimSpace(str)
	}

	var b strings.Builder
	b.Grow(len(str))

	for i = 0; i < len(str); i++ {
		if str[i] == '\\' {
			i++
			if i >= len(str) {
				break
			}
		}
		b.WriteByte(str[i])
	}

	return strings.TrimSpace(b.String())

}
