package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

func handleConnection (connection net.Conn){

	buffer := make ([]byte, 2)
	var final []byte
	var netError net.Error

	defer connection.Close()
	connection.SetReadDeadline(time.Now().Add(100 * time.Second))

	for {
		lettersCount, err := connection.Read(buffer)
		if err != nil {
			if errors.As(err, &netError) && netError.Timeout() {
				fmt.Println("Connection Timed out.")
				} else if err == io.EOF {
					fmt.Println(string(final))
					fmt.Println("End of File")
					final = nil
				}
			break
			}

		for i := 0; i < lettersCount; i++{
			if buffer[i] == '\n' {
					fmt.Println(string(final))
					final = nil
					continue
				}
			final = append(final, buffer[i])
			}
		}
}

func main() {

	listener,_ := net.Listen("tcp", ":9000")

	for {
		fmt.Println("Waiting for a client...")
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener failed to connect")
			break
		}
		fmt.Println("Connection accepted!")
		go handleConnection (connection)
		}
}


