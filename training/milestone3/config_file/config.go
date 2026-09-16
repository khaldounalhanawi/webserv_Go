package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Route struct {
	url		string
	handler	string
}

type Config struct {
	ListenPort		int
	root			string
	routes			[]Route
	error_page		string
	max_header_size	int
	max_body_size	int
	cgi_path		[]string
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
	var needServer		bool
	var config		*Config
	var err			error
	var	tokensLen	int

	tokensLen = len(tokens)

	for i := 0; i < tokensLen; i++ {
 
		// first token must be server
		if i == 0 { needServer = true}
		if needServer && string(tokens[i]) == "server" {
			i ++
			if i >= tokensLen { return nil, errors.New("missing {")}
			needServer = false
			// check for open { after wards
			if string(tokens[i]) == "{" {
				open = true
				i ++
				if i >= tokensLen { return nil, errors.New("missing args")}
				config = new(Config)
			} else {
				return nil, errors.New("Need openning {") }
		} else if needServer {
			return nil, errors.New("Missing Server key word") }

		// switch amongst cases
		switch {
		case string(tokens[i]) == "listen":
			i ++
			if i >= tokensLen { return nil, errors.New("missing listen arg")}
			config.ListenPort, err = strconv.Atoi(string(tokens[i]))
			if err != nil {
				return nil, err }

		case string(tokens[i]) == "root":
			i ++
			if i >= tokensLen { return nil, errors.New("missing root arg")}
			config.root = string(tokens[i])

		case string(tokens[i]) == "}":
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

	var config []*Config
	var tokens [][]byte

	// get tokens from file
	tokens, err := getTokensFromFile(path)
	if err != nil {
		return nil, err }

	// put tokens into config
	config, err = PutTokensInConfig(tokens)
	if err != nil {
		println(err.Error())
		return nil, err }

	// validate configs

	// just for test, print out tokens
	// for n, i := range tokens {
	// 	fmt.Println(n, "_Token:", string(i))
	// }

	for n, i := range config {
		fmt.Println(n, *i)
	}
	return config, nil
}

func main() {
	getConfig("config.config")
}