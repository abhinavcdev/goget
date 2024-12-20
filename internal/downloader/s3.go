// goget/internal/downloader/s3.go
package downloader

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Downloader handles downloading files from S3
type S3Downloader struct {
	client   *s3.Client
	bucket   string
	key      string
	region   string
	filePath string
}

// NewS3Downloader creates a new S3 downloader
func NewS3Downloader(s3URL, outputPath, region string) (*S3Downloader, error) {
	// Parse S3 URL
	bucket, key, err := parseS3URL(s3URL)
	if err != nil {
		return nil, err
	}

	// If output path is a directory, use the key's filename
	fileInfo, err := os.Stat(outputPath)
	if err == nil && fileInfo.IsDir() {
		outputPath = filepath.Join(outputPath, filepath.Base(key))
	} else if outputPath == "" {
		outputPath = filepath.Base(key)
	}

	// Use default region if not specified
	if region == "" {
		region = "us-east-1"
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &S3Downloader{
		client:   client,
		bucket:   bucket,
		key:      key,
		region:   region,
		filePath: outputPath,
	}, nil
}

// Download downloads the file from S3
func (d *S3Downloader) Download() error {
	fmt.Printf("Downloading s3://%s/%s to %s\n", d.bucket, d.key, d.filePath)

	// Create the file
	file, err := os.Create(d.filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %v", d.filePath, err)
	}
	defer file.Close()

	// Get the object from S3
	resp, err := d.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(d.key),
	})
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// Copy the data to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

// parseS3URL parses an S3 URL in the format s3://bucket/key
func parseS3URL(s3URL string) (string, string, error) {
	if !strings.HasPrefix(s3URL, "s3://") {
		return "", "", fmt.Errorf("invalid S3 URL format: %s", s3URL)
	}

	// Parse URL
	u, err := url.Parse(s3URL)
	if err != nil {
		return "", "", err
	}

	bucket := u.Host
	key := strings.TrimPrefix(u.Path, "/")

	if bucket == "" || key == "" {
		return "", "", fmt.Errorf("invalid S3 URL format: %s", s3URL)
	}

	return bucket, key, nil
}
// Updated on December 20, 2024 - Added direct S3 download support
