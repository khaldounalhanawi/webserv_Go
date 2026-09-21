package main

import (
	"loadConfig"
	"parser"
	. "types"
)

func main(){
	var route Route

	// parse & validate Configs
	Configs, err := loadConfig.LoadConfig("config.config")
	if err != nil {
		println(err.Error())
		return }
	println (Configs)

	// get tcp input
	var input []byte // to change

	// parse requests
	// while loop 
	request, _, err := parser.ParseRequest(input)
	if err != nil {
		println(err.Error()) 
		return }
	println(request)

	// for each request
		// request handler
		// give out response
		// tcp
	// loop again



	_,i,err :=  parser.ParseRequest ([]byte{'h','y','a'})
	println (route.Url)
	println (i)
	println (err.Error())
}
