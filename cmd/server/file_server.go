package main

import (
	"net/http"
)

// FileSystem custom file system handler that serves only files and if not found always returns index.html
type fileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs fileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return fs.fs.Open("index.html")
	}

	return f, nil
}

func FrontEndServer(pathToDir string) http.Handler {
	return http.FileServer(fileSystem{http.Dir(pathToDir)})
}
