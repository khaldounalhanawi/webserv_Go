package gethandler

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestGetHandler(t *testing.T) {
	root := t.TempDir()

	// Create files/directories used by the tests.
	err := os.WriteFile(
		filepath.Join(root, "index.html"),
		[]byte("<h1>Hello</h1>"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(
		filepath.Join(root, "hello.txt"),
		[]byte("hello world"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(filepath.Join(root, "images"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(
		filepath.Join(root, "images", "logo.png"),
		[]byte("fake png"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(filepath.Join(root, "empty"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	outside := filepath.Join(filepath.Dir(root), "outside.txt")
	err = os.WriteFile(outside, []byte("secret"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Symlink inside root -> outside root.
	err = os.Symlink(outside, filepath.Join(root, "escape"))
	if err != nil {
		t.Fatal(err)
	}

	config := Config{
		root: root,
	}

	t.Run("GET root returns index.html", func(t *testing.T) {
		request := Request{
			Path: "/",
		}

		result, err := GetHandler(request, Route{}, config)
		if err != nil {
			t.Fatal(err)
		}

		if result.Status != http.StatusOK {
			t.Fatalf("expected 200, got %d", result.Status)
		}

		if string(result.body) != "<h1>Hello</h1>" {
			t.Fatalf("unexpected body: %q", result.body)
		}
	})

	t.Run("GET file", func(t *testing.T) {
		request := Request{
			Path: "/hello.txt",
		}

		result, err := GetHandler(request, Route{}, config)
		if err != nil {
			t.Fatal(err)
		}

		if result.Status != http.StatusOK {
			t.Fatalf("expected 200, got %d", result.Status)
		}

		if string(result.body) != "hello world" {
			t.Fatalf("unexpected body: %q", result.body)
		}

		if result.ContentType != "text/plain; charset=utf-8" {
			t.Fatalf("unexpected content type: %q", result.ContentType)
		}
	})

	t.Run("GET nested file", func(t *testing.T) {
		request := Request{
			Path: "/images/logo.png",
		}

		result, err := GetHandler(request, Route{}, config)
		if err != nil {
			t.Fatal(err)
		}

		if result.Status != http.StatusOK {
			t.Fatalf("expected 200, got %d", result.Status)
		}
	})

	t.Run("GET missing file returns 404", func(t *testing.T) {
		request := Request{
			Path: "/does-not-exist.txt",
		}

		result, err := GetHandler(request, Route{}, config)
		if err != nil {
			t.Fatal(err)
		}

		if result.Status != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", result.Status)
		}
	})

	t.Run("GET directory without index returns 404", func(t *testing.T) {
		request := Request{
			Path: "/empty",
		}

		result, err := GetHandler(request, Route{}, config)
		if err != nil {
			t.Fatal(err)
		}

		if result.Status != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", result.Status)
		}
	})

	t.Run("path traversal is forbidden", func(t *testing.T) {
		request := Request{
			Path: "/../outside.txt",
		}

		result, err := GetHandler(request, Route{}, config)
		if err != nil {
			t.Fatal(err)
		}

		if result.Status != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", result.Status)
		}
	})

	t.Run("symlink outside root is forbidden", func(t *testing.T) {
		request := Request{
			Path: "/escape",
		}

		result, err := GetHandler(request, Route{}, config)
		if err != nil {
			t.Fatal(err)
		}

		if result.Status != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", result.Status)
		}
	})

	t.Run("file with unknown extension has empty MIME type", func(t *testing.T) {
	err := os.WriteFile(
		filepath.Join(root, "README"),
		[]byte("some content"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	request := Request{
		Path: "/README",
	}

	result, err := GetHandler(request, Route{}, config)
	if err != nil {
		t.Fatal(err)
	}

	if result.Status != http.StatusOK {
		t.Fatalf("expected 200, got %d", result.Status)
	}

	if result.ContentType != "" {
		t.Fatalf("expected empty content type, got %q", result.ContentType)
	}
})

t.Run("permission-denied file returns 403", func(t *testing.T) {
	path := filepath.Join(root, "private.txt")

	err := os.WriteFile(
		path,
		[]byte("secret"),
		0000,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Restore permissions when the test finishes so cleanup works.
	t.Cleanup(func() {
		_ = os.Chmod(path, 0644)
	})

	request := Request{
		Path: "/private.txt",
	}

	result, err := GetHandler(request, Route{}, config)
	if err != nil {
		t.Fatal(err)
	}

	// NOTE: this test may fail when running as root,
	// because root can read permission-0000 files.
	if result.Status != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", result.Status)
	}
})

t.Run("symlink inside root is allowed", func(t *testing.T) {
	target := filepath.Join(root, "hello.txt")
	link := filepath.Join(root, "hello-link.txt")

	err := os.Symlink(target, link)
	if err != nil {
		t.Fatal(err)
	}

	request := Request{
		Path: "/hello-link.txt",
	}

	result, err := GetHandler(request, Route{}, config)
	if err != nil {
		t.Fatal(err)
	}

	if result.Status != http.StatusOK {
		t.Fatalf("expected 200, got %d", result.Status)
	}

	if string(result.body) != "hello world" {
		t.Fatalf("unexpected body: %q", result.body)
	}
})

t.Run("directory with index.html is served", func(t *testing.T) {
	dir := filepath.Join(root, "docs")

	err := os.Mkdir(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(
		filepath.Join(dir, "index.html"),
		[]byte("<h1>Docs</h1>"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	request := Request{
		Path: "/docs",
	}

	result, err := GetHandler(request, Route{}, config)
	if err != nil {
		t.Fatal(err)
	}

	if result.Status != http.StatusOK {
		t.Fatalf("expected 200, got %d", result.Status)
	}

	if string(result.body) != "<h1>Docs</h1>" {
		t.Fatalf("unexpected body: %q", result.body)
	}
})

t.Run("root path resolves to root directory", func(t *testing.T) {
	request := Request{
		Path: "/",
	}

	result, err := GetHandler(request, Route{}, config)
	if err != nil {
		t.Fatal(err)
	}

	// Because root contains index.html, "/" should serve it.
	if result.Status != http.StatusOK {
		t.Fatalf("expected 200, got %d", result.Status)
	}

	if string(result.body) != "<h1>Hello</h1>" {
		t.Fatalf("unexpected body: %q", result.body)
	}
})
}