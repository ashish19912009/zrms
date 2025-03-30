package helper

import "time"

func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func NowPtr() *time.Time {
	now := time.Now().UTC()
	return &now
}
