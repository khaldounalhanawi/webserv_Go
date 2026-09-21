package loadConfig

import (
	"errors"
	"strconv"
	. "types"
)

func ParseTokens(tokens [][]byte) ([]ParsedConfig, error) {

	var parsed_configs	[]ParsedConfig
	var open			bool
	var needServer		bool = true
	var config			*Config
	var err				error
	var	tokensLen		int
	var declared		map[string]bool

	tokensLen = len(tokens)
	if tokensLen == 0 {
		return parsed_configs, errors.New("Empty file")}

	for i := 0; i < tokensLen; i++ {
 
		// first token must be server
		if needServer && string(tokens[i]) == "server" {
			if i++; i >= tokensLen { return parsed_configs, errors.New("missing {")}
			needServer = false

			// check for open {
			if string(tokens[i]) == "{" {
				open = true
				if i++; i >= tokensLen { return parsed_configs, errors.New("missing args")}

				// initiate a config
				config = new(Config)
				config.Routes = make([]Route,0)
				declared = make(map[string]bool)
				} else { return parsed_configs, errors.New("Need openning {") }

		} else if needServer {
			return parsed_configs, errors.New("Missing Server key word") }

		// switch amongst cases
		token := string(tokens[i])
		switch token {
		case "listen":
			if declared["listen"] {return parsed_configs, errors.New("Listen Double decalration")}
			if i++; i >= tokensLen { return parsed_configs, errors.New("missing listen arg")}
			config.ListenPort, err = strconv.Atoi(string(tokens[i]))
			if err != nil {
				return parsed_configs, err }
			declared["listen"] = true

		case "root":
			if declared["root"] {return parsed_configs, errors.New("Root Double decalration")}
			if i++; i >= tokensLen { return parsed_configs, errors.New("missing root arg")}
			config.Root = string(tokens[i])
			declared["root"] = true

		case "error_page":
			if declared["error_page"] {return parsed_configs, errors.New("error_page Double decalration")}
			if i++; i >= tokensLen {return parsed_configs, errors.New("missing error_page arg")}
			config.Error_page = string(tokens[i])
			declared["error_page"] = true

		case "route":
			var RouteItem Route
			if i++; i >= tokensLen { return parsed_configs, errors.New("missing route arg")}
			RouteItem.Url = string(tokens[i])
			if i++; i >= tokensLen { return parsed_configs, errors.New("missing route arg")}
			RouteItem.Handler = string(tokens[i])
			config.Routes = append(config.Routes, RouteItem)

		case "max_header_size":
			if declared["max_header_size"] {return parsed_configs, errors.New("max_header_size Double decalration")}
			if i++; i >= tokensLen { return parsed_configs, errors.New("missing max_header_size arg")}
			config.Max_header_size, err = strconv.Atoi(string(tokens[i]))
			if err != nil {
				return parsed_configs, err }
			declared["max_header_size"] = true

		case "max_body_size":
			if declared["max_body_size"] {return parsed_configs, errors.New("max_body_size Double decalration")}
			if i++; i >= tokensLen { return parsed_configs, errors.New("missing max_body_size arg")}
			config.Max_body_size, err = strconv.Atoi(string(tokens[i]))
			if err != nil {
				return parsed_configs, err }
			declared["max_body_size"] = true

		case "}":
			open = false
			needServer = true
			var parsed_config ParsedConfig
			parsed_config.Config = config
			parsed_config.Declared = make(map[string]bool)
			for key, value := range declared {
				parsed_config.Declared[key] = value
			}
			parsed_configs = append(parsed_configs, parsed_config)
			for k := range declared {
				delete(declared, k)
			}
		default :
			return parsed_configs, errors.New("Unknown setting")
		}
	}

	if open {
		return parsed_configs, errors.New("Unclosed bracket") }
	
	return parsed_configs, nil
}
