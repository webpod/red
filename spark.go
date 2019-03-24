package main

import (
	"bytes"
	"math"
)

var steps = []rune("▁▂▃▄▅▆▇") // 8th rune "█" omitted to prevent gluing of rows.

func Spark(nums []float64) string {
	if len(nums) == 0 {
		return ""
	}
	indices := normalize(nums)
	var sparkline bytes.Buffer
	for _, index := range indices {
		sparkline.WriteRune(steps[index])
	}
	return sparkline.String()
}

func normalize(nums []float64) []int {
	var indices []int
	total := float64(len(steps))
	min := minimum(nums)
	for i := range nums {
		nums[i] -= min
	}
	max := maximum(nums)
	if max == 0 {
		// Protect against division by zero
		// This can happen if all values are the same
		max = 1
	}
	for i := range nums {
		x := nums[i]
		x /= max
		x *= total
		if x == total {
			x = total - 1
		} else {
			x = math.Floor(x)
		}
		indices = append(indices, int(x))
	}
	return indices
}

func minimum(nums []float64) float64 {
	var min = nums[0]
	for _, x := range nums {
		if math.Min(x, min) == x {
			min = x
		}
	}
	return min
}

func maximum(nums []float64) float64 {
	var max = nums[0]
	for _, x := range nums {
		if math.Max(x, max) == x {
			max = x
		}
	}
	return max
}
