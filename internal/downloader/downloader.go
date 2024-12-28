// goget/internal/downloader/downloader.go
package downloader

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"
	"time"

	"goget/internal/output"
	"goget/internal/ui"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// Config holds all settings for a download.
type Config struct {
	URL         string
	OutputPath  string
	OutputType  string
	Concurrency int
	Retries     int
	JSONOutput  bool
	S3Bucket    string
	S3Region    string
}

// Downloader manages the file download process.
type Downloader struct {
	config Config
	writer output.TargetWriter
}

// New creates a new Downloader instance.
func New(config Config) (*Downloader, error) {
	// Auto-determine output path if not provided
	if config.OutputPath == "" {
		parsedURL, err := url.Parse(config.URL)
		if err != nil {
			return nil, fmt.Errorf("invalid url: %w", err)
		}
		config.OutputPath = path.Base(parsedURL.Path)
		if config.OutputPath == "" || config.OutputPath == "." {
			config.OutputPath = "index.html"
		}
	}

	var writer output.TargetWriter
	var err error

	// Select the appropriate writer based on the output type
	switch config.OutputType {
	case "file":
		writer, err = output.NewFileWriter(config.OutputPath)
	case "s3":
		writer, err = output.NewS3Writer(config.S3Bucket, config.OutputPath, config.S3Region)
	case "stdout":
		return nil, fmt.Errorf("stdout output not yet implemented")
	default:
		return nil, fmt.Errorf("unknown output type: %s", config.OutputType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create writer: %w", err)
	}

	return &Downloader{config: config, writer: writer}, nil
}

// Download starts the download process.
func (d *Downloader) Download() error {
	defer d.writer.Close()

	// 1. Get file metadata (size, resumability)
	req, err := http.NewRequest("HEAD", d.config.URL, nil)
	if err != nil {
		return err
	}
	
	// Add a User-Agent header to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode > 299 {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	size := int64(-1)
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		size, err = strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			log.Println("Invalid Content-Length. Downloading in a single stream.")
			d.config.Concurrency = 1
		}
	} else {
		// If no content-length, download sequentially
		log.Println("Content-Length not found. Downloading in a single stream.")
		d.config.Concurrency = 1
	}

	// Check for resumable download support
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		log.Println("Server does not support partial downloads. Downloading sequentially.")
		d.config.Concurrency = 1
	}

	// 2. Setup progress tracking
	if d.config.JSONOutput {
		return d.downloadWithJSON(size)
	} else {
		return d.downloadWithProgressBars(size)
	}
}

// downloadWithProgressBars downloads the file with visual progress bars
func (d *Downloader) downloadWithProgressBars(size int64) error {
	// Setup progress bars and concurrency
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))
	
	if size <= 0 {
		// Unknown size, download in a single stream
		wg.Add(1)
		go d.downloadChunk(p, 0, 0, -1, &wg)
	} else {
		// Known size, split into chunks
		chunkSize := int64(math.Ceil(float64(size) / float64(d.config.Concurrency)))
		
		for i := 0; i < d.config.Concurrency; i++ {
			start := int64(i) * chunkSize
			end := start + chunkSize - 1
			if i == d.config.Concurrency-1 {
				end = size - 1
			}
			
			wg.Add(1)
			go d.downloadChunk(p, i, start, end, &wg)
		}
	}
	
	// Wait for all downloads to complete
	wg.Wait()
	
	// 4. Finalize the download
	return d.writer.Commit()
}

// downloadWithJSON downloads the file with JSON progress output
func (d *Downloader) downloadWithJSON(size int64) error {
	// Implementation for JSON output mode
	// This is a simplified version that doesn't use concurrent chunks for JSON output
	req, err := http.NewRequest("GET", d.config.URL, nil)
	if err != nil {
		return err
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode > 299 {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}
	
	// Create JSON progress writer
	jsonProgress := ui.NewJSONProgressWriter(size)
	
	// Copy data through the JSON progress writer
	_, err = io.Copy(io.NewOffsetWriter(d.writer, 0), io.TeeReader(resp.Body, jsonProgress))
	if err != nil {
		return err
	}
	
	return d.writer.Commit()
}

// downloadChunk downloads a specific byte range of the file.
func (d *Downloader) downloadChunk(p *mpb.Progress, id int, start, end int64, wg *sync.WaitGroup) {
	defer wg.Done()

	var err error
	for attempt := 0; attempt <= d.config.Retries; attempt++ {
		err = d.tryDownloadChunk(p, id, start, end)
		if err == nil {
			return // Success!
		}
		log.Printf("Chunk %d failed (attempt %d/%d): %v", id, attempt+1, d.config.Retries, err)
		backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		time.Sleep(backoff)
	}
	log.Fatalf("FATAL: Chunk %d failed after %d retries.", id, d.config.Retries)
}

func (d *Downloader) tryDownloadChunk(p *mpb.Progress, id int, start, end int64) error {
	req, err := http.NewRequest("GET", d.config.URL, nil)
	if err != nil {
		return err
	}
	
	// Add a User-Agent header to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")

	// Only set range header if we're doing a partial request
	if end >= 0 {
		rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
		req.Header.Set("Range", rangeHeader)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if end >= 0 && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server doesn't support partial content: %d", resp.StatusCode)
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var chunkSize int64
	if end >= 0 {
		chunkSize = end - start + 1
	} else {
		// Unknown size
		contentLength := resp.Header.Get("Content-Length")
		if contentLength != "" {
			chunkSize, _ = strconv.ParseInt(contentLength, 10, 64)
		} else {
			chunkSize = 0 // Unknown
		}
	}

	// Create progress bar
	var bar *mpb.Bar
	if chunkSize > 0 {
		bar = ui.NewProgressBar(p, chunkSize, id)
	} else {
		bar = p.AddBar(-1, // Unknown total
			mpb.PrependDecorators(
				decor.Name(fmt.Sprintf("Chunk #%d", id+1), decor.WC{W: 10, C: decor.DindentRight}),
				decor.CountersKibiByte("% .2f / ???"),
			),
			mpb.AppendDecorators(
				decor.EwmaSpeed(decor.SizeB1024(0), "% .2f", 60),
			),
		)
	}

	// Create a proxy reader to update the bar automatically
	proxyReader := bar.ProxyReader(resp.Body)
	defer proxyReader.Close()

	// Copy the data to the correct offset in the output writer
	_, err = io.Copy(io.NewOffsetWriter(d.writer, start), proxyReader)
	return err
}
// Updated on December 3, 2024 - Added core downloader implementation
// Updated on December 28, 2024 - Fixed bugs and improved error handling
