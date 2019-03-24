package main

func ComputeDistance(a, b []string) int {
	if len(a) == 0 {
		return len(b)
	}

	if len(b) == 0 {
		return len(a)
	}

	if equals(a, b) {
		return 0
	}

	// swap to save some memory O(min(a,b)) instead of O(a)
	if len(a) > len(b) {
		a, b = b, a
	}

	// init the row
	x := make([]int, len(a)+1)
	for i := 0; i <= len(a); i++ {
		x[i] = i
	}

	// fill in the rest
	for i := 1; i <= len(b); i++ {
		prev := i
		var current int

		for j := 1; j <= len(a); j++ {

			if b[i-1] == a[j-1] {
				current = x[j-1] // match
			} else {
				current = min(min(x[j-1]+1, prev+1), x[j]+1)
			}
			x[j-1] = prev
			prev = current
		}
		x[len(a)] = prev
	}
	return x[len(a)]
}
