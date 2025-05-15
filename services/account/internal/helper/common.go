package helper

import (
	"time"

	"google.golang.org/grpc/metadata"
)

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

func GetMetadataValue(md metadata.MD, keys ...string) string {
	for _, key := range keys {
		if values := md.Get(key); len(values) > 0 {
			return values[0]
		}
	}
	return ""
}
