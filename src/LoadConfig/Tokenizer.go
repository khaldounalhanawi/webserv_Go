package loadConfig

import (
	"bufio"
	"bytes"
	"os"
)

func TokenizeFile(path string) ([][]byte, error) {

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