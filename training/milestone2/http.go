package main

import (
	"errors"
	"fmt"
	"strings"
)

type Request struct {
	Method	string
	Path	string
	Version	string
	Headers	map[string]string
	Body	[]byte
}

func validateRequestLineTokens(tokens []string) bool {

	if len(tokens) != 3 {
		return false
	}
	if tokens[0] == "" || tokens[1] == ""||
		tokens[2] != "HTTP/1.1" {
		return false
	}
	return true
}

func getRequestLineTokens(dataString string) ([]string, string, error) {

	firstLine, restLines, found := strings.Cut(dataString, "\r\n")
	if !found {
		return nil, restLines, errors.New("Bad format")
	}

	// split first line into tokens
	tokenArray := strings.Split(firstLine, " ")
	return tokenArray, restLines, nil
}

// shouldnt this data be a pointer? so that it is not too heavy copying???
func ParseRequest(data []byte) (Request, error) { 
	var myRequest Request
	var err error
	var dataString string = string(data)

	// if stream empty
	if len(data) == 0 {
		err = errors.New("Empty request")
		return myRequest, err
	}

	// get request line tokens
	requestLineTokens, dataString, err := getRequestLineTokens(dataString)
	if err != nil {
		return myRequest, err
	}

	// validate request line tokens
	if !validateRequestLineTokens(requestLineTokens) {
		err = errors.New("Bad format")
		return myRequest, err
	}

	myRequest.Method = requestLineTokens[0]
	myRequest.Path = requestLineTokens[1]
	myRequest.Version = requestLineTokens[2]

	// find headers
	// go to \n\n and work on before for headers and after for body
	before, dataString, found := strings.Cut(dataString, "\r\n\r\n")
	if !found {
		return myRequest, errors.New("incomplete headers")
	}

	// split before into headers array
	headersArray := strings.Split(before, "\r\n")
	if len(headersArray) == 0 || (len(headersArray) == 1 && headersArray[0] == "") {
		return myRequest, err
	}
	
	// initiate headers array
	myRequest.Headers = make (map[string]string)

	for _, i := range headersArray {
		before, after, found := strings.Cut(i, ":")
		if !found {
			err = errors.New("Bad header format")
			return myRequest, err
		}
		name := strings.ToLower(strings.TrimSpace(before))
		value := strings.TrimSpace(after)
		myRequest.Headers[name] = strings.TrimSpace(value)
	}


	return myRequest, err
}

func PrintRequest(myRequest Request) {

	fmt.Println("Method: ", myRequest.Method)
	fmt.Println("Path: ", myRequest.Path)
	fmt.Println("Version: ", myRequest.Version)
	if len(myRequest.Headers) > 0 {
		for key, val := range myRequest.Headers {
			fmt.Println(key, ": ", val)
		}
	} else {
			fmt.Println("Headers: empty")
	}
	fmt.Println("Body: ", string(myRequest.Body))
}

func main() {
	incoming := []byte("GET /path HTTP/1.1\r\nHost: localhost\r\nUser-Agent: curl\r\nAccept: */*\r\n\r\n")

	myRequest, err := ParseRequest(incoming)
	if err != nil {
		fmt.Println("ERRRRRRRR")
		return
	}
	PrintRequest(myRequest)
}