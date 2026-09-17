package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
)

var MyServerSettings = ServerSettings{
	maximum_header_limit:	8192,
	maximum_body_limit:		1048576,
	default_body_size:		1024,
	default_header_size:	512,
	supports:				[]string{"static", "cgi"},
}

func SplitTokenByBrackets(token []byte) [][]byte {

	var result	[][]byte
	var match	bool
	var last	int = 0

	for n, i := range token {
		if i == '{' || i == '}' {
			if n != last {
				result = append(result, token[last:n]) }
			result = append(result, []byte{i})
			last = n + 1
			match = true }
	}

	if match && last < len(token) {
		result = append (result , token [last:])
	}

	if !match {
		result = append(result, token) }

	return result
}

func getTokensFromFile(path string) ([][]byte, error) {

	var tokens [][]byte

	file, err := os.Open(path)
	if err != nil {
		return nil, err }

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		newTokens := bytes.Fields(scanner.Bytes())
		for _, token := range newTokens {
			if len (token) > 1 {
				splitToken := SplitTokenByBrackets(token)
				for _, i := range splitToken {
					tokens = append(tokens, i)
				}
			} else {
			tokens = append(tokens, token) }
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err }

	return tokens, nil
}

func PutTokensInConfig(tokens [][]byte) ([]*Config, error) {

	var configs 	[]*Config
	var open		bool
	var needServer	bool = true
	var config		*Config
	var err			error
	var	tokensLen	int

	tokensLen = len(tokens)
	if tokensLen == 0 {
		return nil, errors.New("Empty file")}

	for i := 0; i < tokensLen; i++ {
 
		// first token must be server
		if needServer && string(tokens[i]) == "server" {
			if i++; i >= tokensLen { return nil, errors.New("missing {")}
			needServer = false

			// check for open {
			if string(tokens[i]) == "{" {
				open = true
				if i++; i >= tokensLen { return nil, errors.New("missing args")}

				// initiate a config
				config = new(Config)
				config.routes = make([]Route,0)
				config.declared = make(map[string]bool)

				} else { return nil, errors.New("Need openning {") }

		} else if needServer {
			return nil, errors.New("Missing Server key word") }

		// switch amongst cases
		token := string(tokens[i])
		switch token {
		case "listen":
			if config.declared["listen"] {return nil, errors.New("Listen Double decalration")}
			if i++; i >= tokensLen { return nil, errors.New("missing listen arg")}
			config.ListenPort, err = strconv.Atoi(string(tokens[i]))
			if err != nil {
				return nil, err }
			config.declared["listen"] = true

		case "root":
			if config.declared["root"] {return nil, errors.New("Root Double decalration")}
			if i++; i >= tokensLen { return nil, errors.New("missing root arg")}
			config.root = string(tokens[i])
			config.declared["root"] = true

		case "error_page":
			if config.declared["error_page"] {return nil, errors.New("error_page Double decalration")}
			if i++; i >= tokensLen {return nil, errors.New("missing error_page arg")}
			config.error_page = string(tokens[i])
			config.declared["error_page"] = true

		case "route":
			var RouteItem Route
			if i++; i >= tokensLen { return nil, errors.New("missing route arg")}
			RouteItem.url = string(tokens[i])
			if i++; i >= tokensLen { return nil, errors.New("missing route arg")}
			RouteItem.handler = string(tokens[i])
			config.routes = append(config.routes, RouteItem)

		case "max_header_size":
			if config.declared["max_header_size"] {return nil, errors.New("max_header_size Double decalration")}
			if i++; i >= tokensLen { return nil, errors.New("missing max_header_size arg")}
			config.max_header_size, err = strconv.Atoi(string(tokens[i]))
			if err != nil {
				return nil, err }
			config.declared["max_header_size"] = true

		case "max_body_size":
			if config.declared["max_body_size"] {return nil, errors.New("max_body_size Double decalration")}
			if i++; i >= tokensLen { return nil, errors.New("missing max_body_size arg")}
			config.max_body_size, err = strconv.Atoi(string(tokens[i]))
			if err != nil {
				return nil, err }
			config.declared["max_body_size"] = true

		case "}":
			open = false
			needServer = true
			configs = append(configs, config)

		default :
			return nil, errors.New("Unknown setting")
		}
	}

	if open {
		return nil, errors.New("Unclosed bracket") }
	
	return configs, nil
}

func getConfig(path string) ([]*Config, error) {

	var configs []*Config
	var tokens [][]byte

	// get tokens from file
	tokens, err := getTokensFromFile(path)
	if err != nil {
		return nil, err }

	// put tokens into config
	configs, err = PutTokensInConfig(tokens)
	if err != nil {
		println(err.Error())
		return nil, err }

	// validate configs
	err = ValidateConfigs (configs)
	if err != nil {
		println(err.Error())
		return nil, err}

	// just for test, print out tokens
	// for n, i := range tokens {
	// 	fmt.Println(n, "_Token:", string(i))
	// }

	for n, i := range configs {
		fmt.Println(n, *i)
	}
	return configs, nil
}

func main() {
	getConfig("config.config")
}