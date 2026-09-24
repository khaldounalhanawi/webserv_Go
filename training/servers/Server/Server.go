package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

type Config struct {
	Address	string
}

type Server struct {
	Config		Config
	Listener	net.Listener
}

func NewServer(config Config) (*Server, error) {
	
	server := new(Server)
	
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return nil, err }
	
	server.Config = config
	server.Listener = listener

	return server, nil
}

func (s *Server) Start() error {

	fmt.Println("Server listening on", s.Config.Address)

	for {
		connection, err := s.Listener.Accept()
		if err != nil {
			return err }
		go s.HandleConnection(connection)
	}
}

func reader(conn net.Conn) error {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil { return err }
		fmt.Println("client:", string(buffer[:n]))
	}
}

func writer(conn net.Conn) error {
    scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		if _,err := conn.Write([]byte(scanner.Text() + "\n")); err != nil {
			return err }
	}

	return scanner.Err()
}

func (s *Server) HandleConnection(connection net.Conn) error {

	defer connection.Close()

	fmt.Println("Connection established with", connection.RemoteAddr())

	_, err := connection.Write([]byte("Hello hello from my server, write me anything\n"))
	if err != nil {
		return err }

	// buffer := make([]byte, 1024)

	// for {
		// n, err := connection.Read(buffer)
		// if err != nil { return err}
	
		// fmt.Println("client:", string(buffer[:n]))

		// connection.Write([]byte("I GOT THIS:" + string(buffer[:n])))


	// }

	 var wg sync.WaitGroup
	 wg.Add(2)
	 go func() { defer wg.Done(); reader(connection)}()
	 go func() { defer wg.Done(); writer(connection)}()
	 wg.Wait()

	return nil
	}

func (s *Server) CloseServer() error {
	return s.Listener.Close()
}