package infra

import (
	"path/filepath"
	"strings"
)

func trimExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

// max panics when arr does not contain any elements.
func max[T any](arr []T, lessThan func(T, T) bool) T {
	mx := arr[0]
	for _, cur := range arr {
		if lessThan(mx, cur) {
			mx = cur
		}
	}
	return mx
}

const minSeqNo int = 1

func lessThanSeqNo(n1, n2 int) bool {
	return n1 < n2
}

func nextSeqNo(seqNos []int) int {
	if len(seqNos) == 0 {
		return minSeqNo
	}
	return max(seqNos, lessThanSeqNo) + 1
}
