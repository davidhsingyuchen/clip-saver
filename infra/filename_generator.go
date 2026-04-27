package infra

import (
	"path/filepath"
	"strings"
)

func trimExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

// stripTimestampPrefix removes the timestamp prefix (format: YYYY-MM-DD_HH-MM-SS_) from filename.
// Returns the filename unchanged if it doesn't match the expected timestamp format.
func stripTimestampPrefix(filename string) string {
	// Timestamp format: 2006-01-02_15-04-05_ (19 characters for timestamp + 1 for underscore = 20)
	const timestampLength = 19
	const separatorLength = 1
	const totalPrefixLength = timestampLength + separatorLength

	if len(filename) <= totalPrefixLength {
		return filename
	}

	// Check if the prefix matches timestamp pattern (YYYY-MM-DD_HH-MM-SS_)
	// Validate basic structure: chars at positions 4, 7, 10, 13, 16, 19 should be specific characters
	if filename[4] == '-' && filename[7] == '-' && filename[10] == '_' &&
		filename[13] == '-' && filename[16] == '-' && filename[19] == '_' {
		return filename[totalPrefixLength:]
	}

	return filename
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
