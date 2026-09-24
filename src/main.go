package main

import (
	"fmt"
	"loadConfig"
	"os"
	"os/signal"
	"server"
	"syscall"
)

func main(){

	// parse & validate Configs
	configs, err := loadConfig.LoadConfig("config.config")
	if err != nil {
		fmt.Println(err)
		return }

	// create a listening channel for sigterm/ sigint
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	// create an array of server-pointers
	serverArray := make([]*server.Server, 0, len(configs))

	// for each config
	for _, config := range configs {

		// create server
		myServer, err := server.NewServer(*config)
		if err != nil { return } // print it out??

		// append server to list of server for later closing
		serverArray = append(serverArray, myServer)
	}

	// create a channel to report errors back to main
	errChannel := make(chan error, len(serverArray))

	// start servers Concurrently
	for _, srv := range serverArray {
		go func (){
			errChannel <- srv.Start()
		}()
	}

	// wait for the term/interrupt signal or error
	select {
		case sig := <-sigChannel:
			fmt.Println("Received signal ", sig)
		case err := <-errChannel:
			fmt.Println("Error: ", err)
	}

	// close All servers
	for _, srv := range serverArray {
		err := srv.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}
