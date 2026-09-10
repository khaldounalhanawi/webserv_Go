package main

import (
	"bytes"
	"testing"
)

type testCase struct {
	name          string
	data          []byte
	wantErr       bool
	wantMethod    string
	wantPath      string
	wantVersion   string
	wantHost      string
	wantBody      []byte
	wantConsumed  int
}

func TestParseRequest(t *testing.T) {

	tests := []testCase{

		// =========================
		// VALID REQUESTS
		// =========================

		{
			name: "valid GET",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"\r\n",
			),
			wantErr:      false,
			wantMethod:   "GET",
			wantPath:     "/",
			wantVersion:  "HTTP/1.1",
			wantHost:     "localhost",
			wantBody:     nil,
		},

		{
			name: "valid POST",
			data: []byte(
				"POST /test HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"Content-Length: 5\r\n" +
					"\r\n" +
					"hello",
			),
			wantErr:      false,
			wantMethod:   "POST",
			wantPath:     "/test",
			wantVersion:  "HTTP/1.1",
			wantHost:     "localhost",
			wantBody:     []byte("hello"),
		},

		{
			name: "GET with query string",
			data: []byte(
				"GET /search?q=hello HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"\r\n",
			),
			wantErr:      false,
			wantMethod:   "GET",
			wantPath:     "/search?q=hello",
			wantVersion:  "HTTP/1.1",
			wantHost:     "localhost",
			wantBody:     nil,
		},

		{
			name: "POST with empty body",
			data: []byte(
				"POST /submit HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"Content-Length: 0\r\n" +
					"\r\n",
			),
			wantErr:      false,
			wantMethod:   "POST",
			wantPath:     "/submit",
			wantVersion:  "HTTP/1.1",
			wantHost:     "localhost",
			wantBody:     nil,
		},

		{
			name: "header value with spaces",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"User-Agent: my browser\r\n" +
					"\r\n",
			),
			wantErr:      false,
			wantMethod:   "GET",
			wantPath:     "/",
			wantVersion:  "HTTP/1.1",
			wantHost:     "localhost",
			wantBody:     nil,
		},

		{
			name: "lowercase host",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					"host: localhost\r\n" +
					"\r\n",
			),
			wantErr:      false,
			wantMethod:   "GET",
			wantPath:     "/",
			wantVersion:  "HTTP/1.1",
			wantHost:     "localhost",
			wantBody:     nil,
		},

		// =========================
		// BAD REQUEST LINE
		// =========================

		{
			name: "missing HTTP version",
			data: []byte(
				"GET /\r\n" +
					"Host: localhost\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		{
			name: "wrong HTTP version",
			data: []byte(
				"GET / HTTP/2.0\r\n" +
					"Host: localhost\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		{
			name: "too many request line tokens",
			data: []byte(
				"GET / HTTP/1.1 EXTRA\r\n" +
					"Host: localhost\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		{
			name: "empty method",
			data: []byte(
				" / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		{
			name: "empty path",
			data: []byte(
				"GET  HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		// =========================
		// MISSING / INCOMPLETE
		// =========================

		{
			name:    "empty request",
			data:    []byte(""),
			wantErr: true,
		},

		{
			name: "incomplete headers",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					"Host: localhost\r\n",
			),
			wantErr: true,
		},

		{
			name: "no header terminator",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					"Host: localhost",
			),
			wantErr: true,
		},

		{
			name: "incomplete body",
			data: []byte(
				"POST / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"Content-Length: 10\r\n" +
					"\r\n" +
					"hello",
			),
			wantErr: true,
		},

		// =========================
		// BAD HEADERS
		// =========================

		{
			name: "header without colon",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					"Host\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		{
			name: "empty header name",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					": localhost\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		{
			name: "whitespace in header name",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					"Ho st: localhost\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		{
			name: "whitespace before colon",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					"Host : localhost\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		// =========================
		// HOST PROBLEMS
		// =========================

		{
			name: "missing Host",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					"User-Agent: test\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		{
			name: "empty Host",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					"Host:\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		// =========================
		// CONTENT-LENGTH PROBLEMS
		// =========================

		{
			name: "invalid content length",
			data: []byte(
				"POST / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"Content-Length: abc\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		{
			name: "negative content length",
			data: []byte(
				"POST / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"Content-Length: -1\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		{
			name: "plus content length",
			data: []byte(
				"POST / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"Content-Length: +5\r\n" +
					"\r\n" +
					"hello",
			),
			wantErr: true,
		},

		{
			name: "duplicate content length",
			data: []byte(
				"POST / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"Content-Length: 5\r\n" +
					"Content-Length: 5\r\n" +
					"\r\n" +
					"hello",
			),
			wantErr: true,
		},

		// =========================
		// BODY EDGE CASES
		// =========================

		{
			name: "body exactly content length",
			data: []byte(
				"POST / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"Content-Length: 5\r\n" +
					"\r\n" +
					"hello",
			),
			wantErr:     false,
			wantMethod:  "POST",
			wantPath:    "/",
			wantVersion: "HTTP/1.1",
			wantHost:    "localhost",
			wantBody:    []byte("hello"),
		},

		{
			name: "body shorter than content length",
			data: []byte(
				"POST / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"Content-Length: 5\r\n" +
					"\r\n" +
					"hi",
			),
			wantErr: true,
		},

		{
			name: "body longer than content length",
			data: []byte(
				"POST / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"Content-Length: 5\r\n" +
					"\r\n" +
					"helloWORLD",
			),
			wantErr:     false,
			wantMethod:  "POST",
			wantPath:    "/",
			wantVersion: "HTTP/1.1",
			wantHost:    "localhost",
			wantBody:    []byte("hello"),
		},

		// =========================
		// RESOURCE LIMIT ATTACKS
		// =========================

		{
			name: "too many headers",
			data: func() []byte {
				req := "GET / HTTP/1.1\r\nHost: localhost\r\n"

				// Host + 22 additional headers = 23 total
				for i := 0; i < 22; i++ {
					req += "X-Test: value\r\n"
				}

				req += "\r\n"
				return []byte(req)
			}(),
			wantErr: true,
		},

		{
			name: "header value too large",
			data: []byte(
				"GET / HTTP/1.1\r\n" +
					"Host: localhost\r\n" +
					"X-Test: 12345678901234567890123456789012345678901\r\n" +
					"\r\n",
			),
			wantErr: true,
		},

		{
			name: "body too large",
			data: func() []byte {
				bodySize := 10*1024*1024 + 1

				req := []byte(
					"POST / HTTP/1.1\r\n" +
						"Host: localhost\r\n" +
						"Content-Length: 10485761\r\n" +
						"\r\n",
				)

				req = append(req, bytes.Repeat([]byte("A"), bodySize)...)

				return req
			}(),
			wantErr: true,
		},

		{
			name: "total headers too large",
			data: func() []byte {
				hugeValue := bytes.Repeat([]byte("A"), 8000)

				req := []byte(
					"GET / HTTP/1.1\r\n" +
						"Host: localhost\r\n" +
						"X-Test: ",
				)

				req = append(req, hugeValue...)
				req = append(req, []byte("\r\n\r\n")...)

				return req
			}(),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			request, _, err := ParseRequest(tc.data)

			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error state: %v", err)
			}

			if tc.wantErr {
				return
			}

			if request.Method != tc.wantMethod {
				t.Errorf("expected method %q, got %q",
					tc.wantMethod, request.Method)
			}

			if request.Path != tc.wantPath {
				t.Errorf("expected path %q, got %q",
					tc.wantPath, request.Path)
			}

			if request.Version != tc.wantVersion {
				t.Errorf("expected version %q, got %q",
					tc.wantVersion, request.Version)
			}

			if request.Headers["host"] != tc.wantHost {
				t.Errorf("expected host %q, got %q",
					tc.wantHost, request.Headers["host"])
			}

			if string(request.Body) != string(tc.wantBody) {
				t.Errorf("expected body %q, got %q",
					tc.wantBody, request.Body)
			}
		})
	}
}

// =========================
// CONSUMED BYTES
// =========================

func TestParseRequestConsumed(t *testing.T) {

	data := []byte(
		"GET /first HTTP/1.1\r\n" +
			"Host: localhost\r\n" +
			"\r\n" +
			"GET /second HTTP/1.1\r\n" +
			"Host: localhost\r\n" +
			"\r\n",
	)

	request, consumed, err := ParseRequest(data)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if request.Path != "/first" {
		t.Errorf("expected path /first, got %q", request.Path)
	}

	expectedConsumed := len(
		"GET /first HTTP/1.1\r\n" +
			"Host: localhost\r\n" +
			"\r\n",
	)

	if consumed != expectedConsumed {
		t.Errorf(
			"expected consumed %d bytes, got %d",
			expectedConsumed,
			consumed,
		)
	}
}

// =========================
// PARSE TWO REQUESTS
// =========================

func TestParseTwoRequests(t *testing.T) {

	data := []byte(
		"GET /first HTTP/1.1\r\n" +
			"Host: localhost\r\n" +
			"\r\n" +
			"GET /second HTTP/1.1\r\n" +
			"Host: localhost\r\n" +
			"\r\n",
	)

	first, consumed, err := ParseRequest(data)

	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	if first.Path != "/first" {
		t.Fatalf("expected first path /first, got %q", first.Path)
	}

	data = data[consumed:]

	second, consumed, err := ParseRequest(data)

	if err != nil {
		t.Fatalf("second request failed: %v", err)
	}

	if second.Path != "/second" {
		t.Errorf("expected second path /second, got %q", second.Path)
	}

	if consumed != len(data) {
		t.Errorf(
			"expected second request to consume %d bytes, got %d",
			len(data),
			consumed,
		)
	}
}

func TestParseRequestSplitAcrossReads(t *testing.T) {

	part1 := []byte(
		"GET / HTTP/1.1\r\n" +
			"Host: loc",
	)

	part2 := []byte(
		"alhost\r\n" +
			"\r\n",
	)

	// First TCP read: request is incomplete.
	_, consumed, err := ParseRequest(part1)

	if err != ErrIncomplete {
		t.Fatalf("expected ErrIncomplete, got %v", err)
	}

	if consumed != 0 {
		t.Fatalf("expected consumed 0, got %d", consumed)
	}

	// Connection handler receives more bytes and appends them
	// to the existing buffer.
	buffer := append(part1, part2...)

	request, consumed, err := ParseRequest(buffer)

	if err != nil {
		t.Fatalf("unexpected error after receiving more data: %v", err)
	}

	if request.Method != "GET" {
		t.Errorf("expected method GET, got %q", request.Method)
	}

	if request.Path != "/" {
		t.Errorf("expected path /, got %q", request.Path)
	}

	if request.Headers["host"] != "localhost" {
		t.Errorf(
			"expected host localhost, got %q",
			request.Headers["host"],
		)
	}

	if consumed != len(buffer) {
		t.Errorf(
			"expected consumed %d, got %d",
			len(buffer),
			consumed,
		)
	}
}
