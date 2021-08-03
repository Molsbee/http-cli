package service

import (
	"strconv"
	"strings"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

func FormatBytes(bytes int64) string {
	unit := ""
	divisor := 1.0

	switch {
	case bytes >= EXABYTE:
		unit = "EB"
		divisor = EXABYTE
	case bytes >= PETABYTE:
		unit = "PB"
		divisor = PETABYTE
	case bytes >= TERABYTE:
		unit = "TB"
		divisor = TERABYTE
	case bytes >= GIGABYTE:
		unit = "GB"
		divisor = GIGABYTE
	case bytes >= MEGABYTE:
		unit = "MB"
		divisor = MEGABYTE
	case bytes >= KILOBYTE:
		unit = "KB"
		divisor = KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0B"
	case bytes < 0:
		return "Unknown"
	}

	result := strings.TrimSuffix(strconv.FormatFloat(float64(bytes)/divisor, 'f', 2, 64), ".00")
	return result + unit
}
