package parser

import (
	"bytes"
)

func ValidateRequestLineTokens(tokens [][]byte) bool {

	if len(tokens) != 3 {
		return false
	}
	if len(tokens[0]) == 0 || len(tokens[1]) == 0 ||
		!bytes.Equal(tokens[2], []byte("HTTP/1.1")) {
		return false
	}
	return true
}

func GetRequestLineTokens(data []byte) ([][]byte, []byte, error) {

	firstLine, restLines, _ := bytes.Cut(data, []byte("\r\n"))

	// split first line into tokens
	tokenArray := bytes.Split(firstLine, []byte(" "))
	return tokenArray, restLines, nil
}
