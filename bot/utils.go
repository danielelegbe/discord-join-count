package bot

import (
	"database/sql"
	"fmt"
)

func FormatNullIntDuration(minutes sql.NullFloat64) string {
	// Handle nullable values

	if minutes.Valid {
		return fmt.Sprintf("%dh %dm", int(minutes.Float64/60), int(minutes.Float64)%60)
	}

	return "0:00"
}

func FormatDuration(minutes int64) string {
	return fmt.Sprintf("%dh %dm", int(minutes/60), int(minutes)%60)
}
