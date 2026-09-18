package gethandler

import (
	"errors"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GetResult struct {
	Status		int
	ContentType	string
	body		[]byte
}

func GetHandler(request Request, route Route, config Config) (GetResult, error) {

	var getResult GetResult

	// Map url path with root
	joinedPath := filepath.Join(config.root, request.Path)

	// determind if target is a file or directory 
		// directory : look for index.html and read it
		// file : read file
	info, err := os.Stat(joinedPath)
	if err != nil {
		getResult.Status = ResolveError(err)
		return getResult, nil }

	if info.IsDir() {
		joinedPath = filepath.Join(joinedPath, "index.html")}

	// check for Symlinks
	realPath, err := filepath.EvalSymlinks(joinedPath)
	if err != nil {
		getResult.Status = ResolveError(err)
		return getResult, nil}
	realRoot, err := filepath.EvalSymlinks(config.root)
	if err != nil {
		getResult.Status = ResolveError(err)
		return getResult, nil}

	// prevent escaping the root
	outside, err := isOutsideRoot(realRoot, realPath)
	if err != nil {return getResult, err}
	if outside {
		getResult.Status = http.StatusForbidden
		return getResult, nil }

	// read file
	data, err := os.ReadFile(realPath)
	if err != nil { 
		getResult.Status = ResolveError(err)
		return getResult, nil }

	// store data into body
	getResult.body = data

	// determind it MIME type
	getResult.ContentType = mime.TypeByExtension(filepath.Ext(realPath))

	// normal response
	getResult.Status = http.StatusOK

	return getResult, nil
}

func isOutsideRoot (root string, joinedPath string) (bool,error) {

	rel, err := filepath.Rel(root, joinedPath)
	if err != nil { return true, err}
	if rel == ".." || strings.HasPrefix(rel, ".." + string(os.PathSeparator)) {
		return true, nil }
	return false, nil
}

func ResolveError (err error) int {
	if errors.Is(err, os.ErrNotExist) {
		return http.StatusNotFound
	}
	if errors.Is(err, os.ErrPermission) {
		return http.StatusForbidden
	}
	return http.StatusInternalServerError
}
