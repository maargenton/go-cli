package strcase

// FilterParts strips any common prefix fragments from parts that appears in
// groupParts.
func FilterParts(parts, groupParts []string) []string {
	for _, p := range groupParts {
		if len(parts) > 0 && parts[0] == p {
			parts = parts[1:]
		} else {
			break
		}
	}
	return parts
}

// UniqueStrings returns a list of strings with duplicates removed. The input
// list does not have to be sorted, and the order of first appearance is
// preserved in the result.
func UniqueStrings(l []string) []string {
	var m = make(map[string]struct{})
	var r []string
	for _, v := range l {
		if _, ok := m[v]; !ok {
			r = append(r, v)
			m[v] = struct{}{}
		}
	}
	return r
}
