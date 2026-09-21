package parser

import (
	"bytes"
	"errors"
	. "types"
)

const maxBodySize = 10 * 1024 * 1024
const maxHeaderSize = 8 * 1024
var ErrIncomplete = errors.New("incomplete HTTP request")

func ParseRequest(data []byte) (Request, int, error) { 
	var myRequest Request
	myRequest.Headers = make (map[string]string)

	// if stream empty
	if len(data) == 0 {
		return myRequest, 0, ErrIncomplete
	}

	// Find end of headers
	headerEnd := bytes.Index(data, []byte("\r\n\r\n"))
	if headerEnd == -1 {
		return myRequest, 0, ErrIncomplete
	}

	headerData := data[:headerEnd]
	body := data[headerEnd+4:]

	if len(headerData) > maxHeaderSize {
		return myRequest, 0, errors.New("Header too large")
	}

	// get request line tokens
	requestLineTokens, headers, err := GetRequestLineTokens(headerData)
	if err != nil {
		return myRequest, 0, err
	}

	// validate request line tokens
	if !ValidateRequestLineTokens(requestLineTokens) {
		return myRequest, 0, errors.New("Bad format")
	}

	myRequest.Method = string(requestLineTokens[0])
	myRequest.Path = string(requestLineTokens[1])
	myRequest.Version = string(requestLineTokens[2])

	// get headers
	myRequest.Headers, err = GetHeaders(headers)
	if err != nil {
		return myRequest, 0, err
	}

	// validate headers
	err = ValidateHeaders(myRequest.Headers)
	if err != nil {
		return myRequest, 0, err
	}

	// validate content length in header against body
	contentLength, err := GetContentLength(myRequest.Headers)
	if err != nil {
		return myRequest, 0, err
	}
	if contentLength > maxBodySize {
		return myRequest, 0, errors.New("Body too large")
	}
	if len(body) < contentLength {
		myRequest.Body = nil
		return myRequest, 0, ErrIncomplete
	}

	if contentLength == 0 {
		myRequest.Body = nil
	} else {
		myRequest.Body = body[:contentLength]
	}

	consumed := headerEnd + 4 + contentLength

	return myRequest, consumed, nil
}
