package LoadConfig

import (
	"errors"
	"strconv"
)

func ParseTokens(tokens [][]byte) ([]*Config, error) {

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
