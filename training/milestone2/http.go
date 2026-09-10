package main

import (
	"fmt"
)

func main() {
	// incoming := []byte("GET /path HTTP/1.1\r\nHost: localhost\r\nUser-Agent: curl\r\nAccept: */*\r\n\r\nholamola")
	// incoming := []byte("POST / HTTP/1.1\r\nContent-Length: 5\r\n\r\nhello")
	incoming := []byte("POST / HTTP/1.1\r\nContent-Length: 5\r\nhost: BMW\r\n\r\nhello")
	// incoming := []byte("POST / HTTP/1.1\r\nContent-Length: 5\r\ncontent-length: 9\r\n\r\nhello")
	// incoming := []byte("POST / HTTP/1.1\r\nContent-Length: 5\r\n\r\nhel")
	// incoming := []byte("POST / HTTP/1.1\r\nContent-Length: 5\r\n\r\nhelloWORLD")
	// incoming := []byte("POST / HTTP/1.1\r\nhost: hola\r\nContent-Length: 5\r\n\r\nhelloWORLD")
	// incoming := []byte("POST / HTTP/1.1\r\nhost: hola\r\nho@m: bis\r\nContent-Length: 5\r\n\r\nhelloWORLD")
	// incoming := []byte("POST /path HTTP/1.1\r\n\r\n")

	myRequest, err := ParseRequest(incoming)
	if err != nil {
		fmt.Println(err)
		return
	}
	PrintRequest(myRequest)
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