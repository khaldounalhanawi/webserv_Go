package LoadConfig

import (
	"fmt"
)

func LoadConfig(path string) ([]*Config, error) {

	var configs []*Config
	var tokens [][]byte

	// get tokens from file
	tokens, err := TokenizeFile(path)
	if err != nil {
		return nil, err }

	// put tokens into config
	configs, err = ParseTokens(tokens)
	if err != nil {
		println(err.Error())
		return nil, err }

	// validate configs
	err = ValidateConfigs(configs)
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
