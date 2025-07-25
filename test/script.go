package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	bucket := "skillflix-bucket"
	key := "evault-demo.mp4"
	localPath := filepath.Base(key)

	start := time.Now()

	err := DownloadFromS3(bucket, key, localPath)
	if err != nil {
		log.Fatalf("Failed to download video: %v", err)
	}

	elapsed := time.Since(start)

	fmt.Println("Downloaded:", localPath)
	fmt.Printf("⏱️ Time taken: %.2f seconds\n", elapsed.Seconds())
}
func DownloadFromS3(bucket, key, localPath string) error {
	ctx := context.Background()

	// Load AWS credentials and config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-south-1")) // Replace with your actual region

	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Create GetObject request
	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer resp.Body.Close()

	// Create local file
	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer outFile.Close()

	// Copy content from S3 to file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file locally: %w", err)
	}

	return nil
}
