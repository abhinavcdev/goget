// goget/internal/output/writer.go
package output

import "io"

// TargetWriter defines the interface for any output destination.
type TargetWriter interface {
	io.WriterAt
	io.Closer
	// Commit finalizes the download (e.g., renaming a .part file)
	Commit() error
}
