package domain

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const (
	fileNameSeparator = "_"
	fileExt           = "png"
)

// FullFileName returns the full file name by
// prepending a timestamp and appending a file extension.
func FullFileName(fileName string) string {
	timestamp := time.Now().Format("2006-01-02" + fileNameSeparator + "15-04-05")
	return fmt.Sprintf("%s%s%s.%s", timestamp, fileNameSeparator, fileName, fileExt)
}

// ExtractFileName returns the file name after
// stripping the timestamp prefix and the file extension suffix.
func ExtractFileName(rawFileName string) string {
	parts := strings.Split(rawFileName, fileNameSeparator)
	return trimExtension(parts[len(parts)-1])
}

func trimExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
