// goget/cmd/root.go
package cmd

import (
	"fmt"
	"goget/internal/downloader"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	output      string
	concurrency int
	retries     int
	jsonOutput  bool
	outputType  string // "file", "s3", "stdout"
	s3Bucket    string
	s3Region    string
)

var rootCmd = &cobra.Command{
	Use:   "goget [URL]",
	Short: "goget is a faster, concurrent file downloader built in Go.",
	Long:  `A modern replacement for wget with support for concurrency, resumable downloads, JSON output, and pluggable storage backends like S3.`,
	Args:  cobra.ExactArgs(1), // Ensures exactly one URL is passed
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		// Check if it's an S3 URL
		if strings.HasPrefix(url, "s3://") {
			// Use S3 downloader directly
			s3dl, err := downloader.NewS3Downloader(url, output, s3Region)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating S3 downloader: %v\n", err)
				os.Exit(1)
			}

			if err := s3dl.Download(); err != nil {
				fmt.Fprintf(os.Stderr, "S3 download failed: %v\n", err)
				os.Exit(1)
			}

			if !jsonOutput {
				fmt.Println("\n✅ Download complete.")
			} else {
				fmt.Println("{\"status\": \"complete\"}")
			}
			return
		}

		// Regular HTTP download
		config := downloader.Config{
			URL:         url,
			Concurrency: concurrency,
			Retries:     retries,
			JSONOutput:  jsonOutput,
			OutputType:  outputType,
			OutputPath:  output,
			S3Bucket:    s3Bucket,
			S3Region:    s3Region,
		}

		// Create and run the downloader
		dl, err := downloader.New(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating downloader: %v\n", err)
			os.Exit(1)
		}

		if err := dl.Download(); err != nil {
			fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
			os.Exit(1)
		}

		if !jsonOutput {
			fmt.Println("\n✅ Download complete.")
		} else {
			fmt.Println("{\"status\": \"complete\"}")
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output file name (defaults to filename from URL)")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 8, "Number of concurrent download connections")
	rootCmd.Flags().IntVarP(&retries, "retries", "r", 5, "Number of retries on failed chunk download")
	rootCmd.Flags().BoolVar(&jsonOutput, "json", false, "Enable machine-readable JSON output")
	rootCmd.Flags().StringVar(&outputType, "output-type", "file", "Output type: file, s3, or stdout")
	rootCmd.Flags().StringVar(&s3Bucket, "s3-bucket", "", "S3 bucket name (required when output-type is s3)")
	rootCmd.Flags().StringVar(&s3Region, "s3-region", "us-east-1", "S3 region (used when output-type is s3)")
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
