package main

import (
	"loadConfig"
	"server"
)

func main(){

	// parse & validate Configs
	configs, err := loadConfig.LoadConfig("config.config")
	if err != nil {
		println(err.Error())
		return }

	// for each config
	for _, config := range configs {

		// create server
		myServer, err := server.NewServer(*config)
		if err != nil { return } // print it out??
		
		// start server
		err = myServer.Start()
		if err != nil { return }

		// close server
		err = myServer.Close()
		if err != nil { return }
	}

	return
}
