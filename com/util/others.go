package util

func SplitRule(c rune) bool {
	if c == '\t' || c == ' ' {
		return true
	} else {
		return false
	}
}
