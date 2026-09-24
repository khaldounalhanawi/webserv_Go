package server

import (
	"net"
	"strconv"
	. "types"
)

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

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
		return err }

		go s.HandleConnection(conn)
	}
}

func (s *Server) Close() error {
	return s.Listener.Close()
}

func (s *Server) HandleConnection(connection net.Conn) error {

	defer connection.Close()

	// parse request

	// request handler

	// serialize response

	// conn.Write(response)

	return nil // temp
}