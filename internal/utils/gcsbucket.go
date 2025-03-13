package utils

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"time"
)

func UploadFileToGCS(file io.Reader, fileName, bucketName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create GCS client: %v", err)
		return "", err
	}
	defer client.Close()

	// Generate a unique filename (optional)
	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileName)

	// Define GCS object path
	objectPath := fmt.Sprintf("uploads/%s", uniqueFileName)

	// Create object handle
	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectPath)
	writer := object.NewWriter(ctx)

	// Copy file content to GCS
	if _, err := io.Copy(writer, file); err != nil {
		log.Printf("Failed to upload file to GCS: %v", err)
		return "", err
	}

	// Close writer to complete upload
	if err := writer.Close(); err != nil {
		log.Printf("Failed to finalize upload: %v", err)
		return "", err
	}

	// Return the full GCS path
	return fmt.Sprintf("gs://%s/%s", bucketName, objectPath), nil
}
func DeleteFileFromGCS(gcsURL, bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create GCS client: %v", err)
		return err
	}
	defer client.Close()
	filePath := GetFilePathFromGCSURL(gcsURL)
	// Create object handle
	bucket := client.Bucket(bucketName)
	object := bucket.Object(filePath)
	log.Printf("Deleting file: gs://%s/%s", bucketName, filePath)
	// Delete the file
	if err := object.Delete(ctx); err != nil {
		log.Printf("Failed to delete file from GCS: %v", err)
		return err
	}

	log.Printf("Successfully deleted file: gs://%s/%s", bucketName, filePath)
	return nil
}

// write a function to trim filename out of filepath: gs://shaw-bucket/uploads/1741664748290248000-image.png i want to get /uploads/1741664748290248000-image.png
func GetFilePathFromGCSURL(gcsURL string) string {
	parts := strings.SplitN(gcsURL, "/", 4)
	if len(parts) < 4 {
		return ""
	}
	return parts[3]
}
