package main

import "sort"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func equals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}
