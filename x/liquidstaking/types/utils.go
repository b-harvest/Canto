package types

import "sort"

func FilterSlice[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func MapToSortedSlice[T any](m map[uint64]T, f func(i, j T) bool) []T {
	var r []T
	for _, v := range m {
		r = append(r, v)
	}
	sort.Slice(r, func(i, j int) bool {
		return f(r[i], r[j])
	})
	return r
}
