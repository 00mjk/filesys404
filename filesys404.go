// Copyright (c) 2021 Abhijit Bose. All Right reserved.
// Use of this source code is governed by a Apache 2.0 license that can be found
// in the LICENSE file.

/*
Package filesys404 helps to implement custom 404 request for FileSystem Queries.

Features of this package:
 - Standard plugging using `FileSystem` type of `net/http` package
 - Fully compatible with `DefaultServeMux` of `net/http` package
 - Protect `.dot` files or hidden files from being served
 - Redirect the not found request to a pre-define custom `Handler`
 - Can be used with custom routers like
   - https://github.com/julienschmidt/httprouter
   - https://github.com/go-chi/chi

Here is an example of how this library can be used:

	dir := http.Dir("./static") // Static Assets directory
	fs := filesys404.New(dir, func(w http.ResponseWriter, r *http.Request){
		// Custom 404 Message
		http.Error(w, "We are sorry - this resource was not found", http.StatusNotImplemented)
	})
	http.HandleFunc("/static",http.StripPrefix("/static", fs))

To Install this package:

	go get -u github.com/boseji/filesys404

Docs available at : https://pkg.go.dev/github.com/boseji/filesys404
*/
package filesys404

import (
	"net/http"
	"path"
	"strings"
)

// FileSystemWith404 stores the supplied static file system and
// the custom Not found handler.
type FileSystemWith404 struct {
	root     http.FileSystem
	notFound http.HandlerFunc
}

// New creates a new FileSystem404 instance
func New(r http.FileSystem, notFound http.HandlerFunc) *FileSystemWith404 {
	return &FileSystemWith404{
		root:     r,
		notFound: notFound,
	}
}

// ServeHTTP is the implementation of the Handler interface
func (fs *FileSystemWith404) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const indexPage = "/index.html"
	// Find out the Path
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	upath = path.Clean(upath)

	// Filter out .files or hidden dot files
	dotFound := false
	for _, p := range strings.Split(r.URL.Path, "/")[1:] {
		if strings.HasPrefix(p, ".") {
			dotFound = true
			break
		}
	}
	if dotFound {
		fs.notFound(w, r)
		return
	}

	// Replace or Dir Lising to Index Pages
	if strings.HasSuffix(r.URL.Path, "/") {
		upath = path.Join(r.URL.Path, indexPage)
	}

	// Try to Open the File
	f, err := fs.root.Open(upath)
	if err != nil {
		// Else its actually an Invalid file
		fs.notFound(w, r)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		fs.notFound(w, r)
		return
	}

	if d.IsDir() {
		f.Close() // Force Close the Directory

		// Check if its just a Dir name that might contain an Index file
		url := r.URL.Path
		if url[len(url)-1] != '/' { // Does not have a '/' at the end
			p := path.Base(url) + "/"
			localRedirect(w, r, p)
			return
		}

		// For Suppressing Directory Listing
		fs.notFound(w, r)
		return
	}

	// Serve the file since we know it actually exists
	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

// localRedirect gives a Moved Permanently response.
// It does not convert relative paths to absolute paths like Redirect does.
func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}
