package tools

import "strings"

func SliceHasPrefix(prifixes []string, path string) bool {
	for _, prefix := range prifixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
