package contentorigin

import "strings"

// ByAlpha sorts paths alphabetically
type ByAlpha []string

func (s ByAlpha) Len() int {
	return len(s)
}
func (s ByAlpha) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByAlpha) Less(i, j int) bool {
	ilower := strings.ToLower(s[i])
	jlower := strings.ToLower(s[j])
	if ilower == jlower {
		return s[i] < s[j]
	}
	return ilower < jlower
}
