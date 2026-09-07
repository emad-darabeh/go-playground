package main

import "log/slog"

func main() {
	arr1 := []int{1, 2, 3, 3, 4, 5}
	arr2 := []int{3, 3, 5, 7}

	result := intersection(arr1, arr2)
	slog.Info("result", "res", result)
}

// two sorted arrays
// [1, 2, 3, 4, 5], [3, 5, 7]

// the output of the intersection: [3, 5]
func intersection(arr1, arr2 []int) []int {
	var result []int

	i1, i2 := 0, 0

	for i1 < len(arr1) && i2 < len(arr2) {
		if arr1[i1] == arr2[i2] {
			//slog.Info("test", "len of result", len(result))
			if len(result) == 0 || arr1[i1] != result[len(result)-1] {
				result = append(result, arr1[i1])
			}
			i1++
			i2++
			continue
		}

		if arr1[i1] < arr2[i2] {
			i1++
		} else {
			i2++
		}
	}

	return result
}
