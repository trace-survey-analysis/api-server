package utils

import (
	"strings"
)

// Extracts the bucket name from a GCS URL
func ExtractBucketNameFromGCS(gcsURL string) string {
	// Format: gs://bucket-name/path/to/file
	parts := strings.Split(gcsURL, "/")
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}

// Extracts the file path from a GCS URL
func ExtractFilePathFromGCS(gcsURL string) string {
	// Format: gs://bucket-name/path/to/file
	parts := strings.Split(gcsURL, "/")
	if len(parts) < 4 {
		return ""
	}
	return strings.Join(parts[3:], "/")
}
