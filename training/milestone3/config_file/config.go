package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

type Config struct {

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


func getConfig(path string) (*Config, error) {

	var config *Config
	var tokens [][]byte

	// get tokens from file
	tokens, err := getTokensFromFile(path)
	if err != nil {
		return nil, err }



	// just for test, print out tokens
	for n, i := range tokens {
		fmt.Println(n, "_Token:", string(i))
	}

	return config, nil
}

func main() {
	getConfig("config.config")
}