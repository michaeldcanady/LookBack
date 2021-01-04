package conversion

import (
	"fmt"
)

func ByteCountSI(b int64, unit int64, pad int) string {

	if b < unit {
		return fmt.Sprintf("%*d B", pad, b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	// using * includes decimal places
	return fmt.Sprintf("%0*.1f %cB", pad, (float64(b) / float64(div)), "kMGTPE"[exp])
}

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
