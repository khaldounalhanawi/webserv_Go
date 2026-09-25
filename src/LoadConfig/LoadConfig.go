package loadConfig

import (
	"fmt"
	. "types"
)

func LoadConfig(path string) ([]*Config, error) {

	var tokens [][]byte

	tokens , err := TokenizeFile(path)
	if err != nil {
		return nil, fmt.Errorf("@LoadConfig.TokenizeFile: %v", err) }

	parsedConfigs, err := ParseTokens(tokens)
	if err != nil {
		return nil, fmt.Errorf("@LoadConfig.ParseTokens: %v", err) }

	err = ValidateConfigs(parsedConfigs)
	if err != nil {
		return nil, fmt.Errorf("@LoadConfig.ValidateConfigs: %v", err) }

	finalConfigsArray := make([]*Config, 0)
	for _, parsedConfig := range parsedConfigs {
		finalConfigsArray = append(finalConfigsArray, parsedConfig.Config)
	}

	return finalConfigsArray, nil
}
