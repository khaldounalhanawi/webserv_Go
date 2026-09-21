package loadConfig

import (
	. "types"
)

func LoadConfig(path string) ([]*Config, error) {

	var tokens [][]byte

	tokens , err := TokenizeFile(path)
	if err != nil {
		return nil, err }

	parsedConfigs, err := ParseTokens(tokens)
	if err != nil {
		println(err.Error())
		return nil, err }

	err = ValidateConfigs(parsedConfigs)
	if err != nil {
		println(err.Error())
		return nil, err }

	return nil, nil
}
