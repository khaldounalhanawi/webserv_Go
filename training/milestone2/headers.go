package main

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

func GetHeaders(headers []byte) (map[string]string, error) {

	headersCountLimit := 22
	headerCharsLimit := 40
	headersMap := make(map[string]string)

	// for empty headers
	if len(headers) == 0 {
		return headersMap, nil
	}

	// split headers by \r\n
	headersArray := bytes.Split(headers, []byte("\r\n"))
	if len(headersArray) == 0 || (len(headersArray) == 1 && len(headersArray[0]) == 0) {
		return headersMap, errors.New("bad headers")
	}

	// header count limit
	if len(headersArray) > headersCountLimit {
		return headersMap, errors.New("Too many headers")
	}

	// split header elements inside array by ':'
	for _, i := range headersArray {

		before, after, found := bytes.Cut(i, []byte(":"))
		if !found {
			return headersMap, errors.New("bad headers")
			}
		// limit size of header
		if len(before) > headerCharsLimit || len(after) > headerCharsLimit {
			return headersMap, errors.New("Too many characters in header")
		}

		// prevent empty name
		name := string(bytes.ToLower(before))
		if len(name) == 0 {
			return headersMap, errors.New("bad headers")
			}
		// only allowed characters in header name
		for _, c := range before {
			if !isTChar(c) {
				return headersMap, errors.New("Invalid character in header name")
				}
			}

		value := bytes.TrimSpace(after)

		// reject content length duplication
		if _, exists := headersMap[name]; exists && name == "content-length" {
			return headersMap, errors.New("Content-length duplicate")
			}

		headersMap[name] = string(value)
		}

	return headersMap, nil
}

func	GetContentLength(headers map[string]string) (int, error) {

	value, exists := headers["content-length"]
	if !exists {
		return 0, nil
	}
	if strings.Contains(value, "+") {
		return 0, errors.New("Bad content-length value")
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	} else if intVal < 0 {
		return 0, errors.New("Bad content-length value")
	}
	return intVal, nil
}


func	ValidateHeaders(headers map[string]string) error {

	// No host or empty host
	if _, exists := headers["host"]; !exists {
		return errors.New("No Host exists")
	} else if headers["host"] == "" {
		return errors.New("Empty host")
	}
	return nil
}

func	isTChar(c byte) bool {

	switch c {
	case '!', '#', '$', '%', '&', '\'', '*',
		'+', '-', '.', '^', '_', '`', '|', '~':
		return true
	}

	// ALPHA
	if (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') {
		return true
	}

	// DIGIT
	if c >= '0' && c <= '9' {
		return true
	}

	return false
}