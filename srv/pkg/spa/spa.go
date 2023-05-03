package spa

import (
	"net/http"
	"os"
)

const (
	IndexFile = "index.html"
)

// FileSystem wraps a standard http.FileSystem however
// it intercepts 404 errors and returns the default
// index.html page instead.
type FileSystem struct {
	root http.FileSystem
}

func NewFileSystem(root http.FileSystem) *FileSystem {
	return &FileSystem{
		root: root,
	}
}

func (fs *FileSystem) Open(name string) (http.File, error) {
	f, err := fs.root.Open(name)
	if os.IsNotExist(err) {
		return fs.root.Open(IndexFile)
	}
	return f, err
}
