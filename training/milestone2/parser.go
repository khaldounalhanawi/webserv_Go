package main

import (
	"bytes"
	"errors"
)

const maxBodySize = 10 * 1024 * 1024 // 10 MiB
const maxHeaderSize = 10 * 1024 * 1024 // 10 MiB

func ParseRequest(data []byte) (Request, error) { 
	var myRequest Request
	var err error

	// if stream empty
	if len(data) == 0 {
		err = errors.New("Empty request")
		return myRequest, err
	}

	// initiate headers array
	myRequest.Headers = make (map[string]string)

	// Find end of headers
	data, body, found := bytes.Cut(data, []byte("\r\n\r\n"))
	if !found {
		return myRequest, errors.New("incomplete headers")
	}

	if len(data) > maxHeaderSize {
		return myRequest, errors.New("Header too large")
	}
	if len(body) > maxBodySize {
		return myRequest, errors.New("Body too large")
	}

	// get request line tokens
	requestLineTokens, headers, err := GetRequestLineTokens(data)
	if err != nil {
		return myRequest, err
	}

	// validate request line tokens
	if !ValidateRequestLineTokens(requestLineTokens) {
		err = errors.New("Bad format")
		return myRequest, err
	}

	myRequest.Method = string(requestLineTokens[0])
	myRequest.Path = string(requestLineTokens[1])
	myRequest.Version = string(requestLineTokens[2])

	// get headers
	myRequest.Headers, err = GetHeaders(headers)
	if err != nil {
		return myRequest, err
	}

	// validate headers
	err = ValidateHeaders(myRequest.Headers)
	if err != nil {
		return myRequest, err
	}

	// validate headers against body
	length, err := GetContentLength(myRequest.Headers)
	if err != nil {
		return myRequest, err
	}
	if length == 0 {
		myRequest.Body = nil
	}
	if length == 0 && len(body) != 0 {
		return myRequest, errors.New("400 Bad Request")
	}
	if len(body) < length {
		myRequest.Body = nil
		return myRequest, errors.New("incomplete body")
	}
	if length > maxBodySize {
		return myRequest, errors.New("Body too large")
	}

	// add body
	myRequest.Body = body[:length]
	return myRequest, err
}

// // getContentLength TESTER
// func main() {

// 	var tester map[string]string =  map[string]string{"banana" : "one","content-length":"5"}

// 	val, err := getContentLength (tester)
// 	if err != nil {
// 		fmt.Println (err)
// 		return
// 	}
// 	fmt.Println(val)
// }

