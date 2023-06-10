package contentorigin

import (
	"os"
	"time"
)

// FileInfo implements os.FileInfo for remote files
type FileInfo struct {
	name string
}

// Name ...
func (f *FileInfo) Name() string {
	panic("not implemented") // TODO: Implement
}

// Size ...
func (f *FileInfo) Size() int64 {
	panic("not implemented") // TODO: Implement
}

// Mode ...
func (f *FileInfo) Mode() os.FileMode {
	panic("not implemented") // TODO: Implement
}

// ModTime ...
func (f *FileInfo) ModTime() time.Time {
	panic("not implemented") // TODO: Implement
}

// IsDir ...
func (f *FileInfo) IsDir() bool {
	panic("not implemented") // TODO: Implement
}

// Sys ...
func (f *FileInfo) Sys() interface{} {
	panic("not implemented") // TODO: Implement
}
