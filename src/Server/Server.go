package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"parser"
	"strconv"
	"time"
	. "types"
)

// Set time out value for reading from servers
const readTimeOut time.Duration = time.Second * 10

type Server struct {
	Config		Config
	Listener	net.Listener
}

func NewServer(config Config) (*Server, error) {

	var myServer *Server = new(Server)
	var err error

	myServer.Config = config
	myServer.Listener, err = net.Listen(
							"tcp",
							":" + strconv.Itoa(config.ListenPort),
							)
	return myServer, err
}

func (s *Server) Start() error {

	fmt.Printf("Server listening on port :%d ..", s.Config.ListenPort)

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
		return err }
		
		go func(){
			err := s.HandleConnection(conn)
			if err != nil {
				fmt.Println(err) }
		}()
	}

}

func (s *Server) Close() error {
	return s.Listener.Close()
}

func (s *Server) HandleConnection(connection net.Conn) error {

	defer connection.Close()

	fmt.Println("Connection: listening @", s.Config.ListenPort)

	leftover := make([]byte,0)

	for {
		// add time out
		connection.SetReadDeadline(time.Now().Add(readTimeOut))

		// get request
		request, newLeftover, err := getRequest(connection, leftover)
		if err != nil {
			if errors.Is(io.EOF, err) {
				return nil }
			return err }

		// request handler

		// serialize response

		// conn.Write(response)

		leftover = newLeftover

		fmt.Println(request) // temp
		// _=request// temp
	}
}

func getRequest(connection net.Conn, leftover []byte) (*Request, []byte, error) {

	buffer := make([]byte, 4096)
	data := make([]byte, 0)

	// add left over data from previous calls
	data = append(data, leftover...)

	for {
		// read from connection
		n, err := connection.Read(buffer)
		if err != nil {
			return nil, nil, err }
		// add read bytes to data
		data = append(data, buffer[:n]...)
		// get request and number of bytes used
		request, used, err := parser.ParseRequest(data)
		switch err {
		// request is complete =
		case nil:
			return &request, data[used:], nil
		// incomplete request =
		case parser.ErrIncomplete:
			continue
		// some error happened need to report
		default:
			return nil, nil, err
		}
	}
}