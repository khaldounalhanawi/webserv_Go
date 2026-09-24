package main

import (
	"server"
)

var myConfig server.Config = server.Config {Address:":8080"}

func main(){
	serv, err := server.NewServer(myConfig)
	if err != nil {
		println(err.Error())
		return }

	serv.Start()
	serv.CloseServer()
}