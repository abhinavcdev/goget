// goget/internal/output/file.go
package output

import (
	"os"
)

// FileWriter implements TargetWriter for local file output
type FileWriter struct {
	file      *os.File
	finalPath string
	partPath  string
}

// NewFileWriter creates a writer for a local file.
func NewFileWriter(path string) (*FileWriter, error) {
	partPath := path + ".gogetpart"

	f, err := os.OpenFile(partPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		file:      f,
		finalPath: path,
		partPath:  partPath,
	}, nil
}

func (fw *FileWriter) WriteAt(p []byte, off int64) (n int, err error) {
	return fw.file.WriteAt(p, off)
}

func (fw *FileWriter) Close() error {
	return fw.file.Close()
}

// Commit renames the partial file to its final destination name.
func (fw *FileWriter) Commit() error {
	// Ensure file is closed before renaming
	fw.Close()
	return os.Rename(fw.partPath, fw.finalPath)
}
