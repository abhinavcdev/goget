// goget/internal/ui/progress.go
package ui

import (
	"fmt"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// NewProgressBar creates a new progress bar for a download chunk
func NewProgressBar(p *mpb.Progress, total int64, id int) *mpb.Bar {
	name := fmt.Sprintf("Chunk #%d", id+1)
	return p.AddBar(total,
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DindentRight}),
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.SizeB1024(0), "% .2f", 60),
		),
	)
}

// JSONProgressWriter implements a progress writer that outputs JSON
type JSONProgressWriter struct {
	totalBytes    int64
	downloadedBytes int64
}

// NewJSONProgressWriter creates a new JSON progress writer
func NewJSONProgressWriter(total int64) *JSONProgressWriter {
	return &JSONProgressWriter{
		totalBytes: total,
	}
}

// Write implements io.Writer and updates the progress
func (j *JSONProgressWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	j.downloadedBytes += int64(n)
	j.printProgress()
	return n, nil
}

// printProgress outputs the current progress as JSON
func (j *JSONProgressWriter) printProgress() {
	percentage := float64(j.downloadedBytes) / float64(j.totalBytes) * 100
	fmt.Printf("{\"progress\": %.2f, \"downloaded\": %d, \"total\": %d}\n", 
		percentage, j.downloadedBytes, j.totalBytes)
}
// Updated on December 5, 2024 - Added progress bar UI
