package loadConfig

import (
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
	. "types"
)

func ValidateConfigs(parsedConfigs []ParsedConfig) error {

	usedPorts := make ([]int,0)

	for n, parsedConfig := range parsedConfigs {

		config := parsedConfig.Config

		// must have (listen, root, at least one route)
		if !parsedConfig.Declared["listen"] {return fmt.Errorf("Server %d: Missing Listen in Config file", n)}
		if !parsedConfig.Declared["root"] {return fmt.Errorf("Server %d: Missing root in Config file", n)}
		if len(config.Routes) < 1 {return fmt.Errorf("Server %d: Must have atleast one Route", n)}

		// port: must be between 1 and 65535
		if config.ListenPort < 1 || config.ListenPort > 65535 {
			return fmt.Errorf("Server %d: Not a valid port number", n)}
		// port: must not use already used port
		if !slices.Contains(usedPorts, config.ListenPort) {
			usedPorts = append(usedPorts, config.ListenPort)
			} else {
				return fmt.Errorf("Server %d: Port %d is already in use", n, config.ListenPort)}

		// root: The root path must not be empty
		if config.Root == "" {return fmt.Errorf("Server %d: Root path is empty", n)}
		// root: The root path must refer to an existing directory
		dir, err := os.Stat(config.Root)
		if err != nil {return fmt.Errorf("Server %d: %s", n, err.Error())}
		if !dir.IsDir() {return fmt.Errorf("Server %d: Root is not a Directory", n)}

		if parsedConfig.Declared["error_page"] {
			err := ValidateErrorPage(config.Error_page)
			if err != nil {
				return fmt.Errorf("Server %d: %s", n, err.Error())}
		}

		// Max header size must be between 0 and limit from server settings
		if parsedConfig.Declared["max_header_size"] {
			if config.Max_header_size < 0 || config.Max_header_size > MyServerSettings.Maximum_header_limit {
				return fmt.Errorf("Server %d: Max header size must be between 0 and %d", n, MyServerSettings.Maximum_header_limit)}
		} else {
			config.Max_header_size = MyServerSettings.Default_header_size }

		// Max body size must be between 0 and limit from server settings
		if parsedConfig.Declared["max_body_size"] {
			if config.Max_body_size < 0 || config.Max_body_size > MyServerSettings.Maximum_body_limit {
				return fmt.Errorf("Server %d: Max body size must be between 0 and %d", n, MyServerSettings.Maximum_body_limit)}
		} else {
			config.Max_body_size = MyServerSettings.Default_body_size }
		
		// Routes: validate routes
		err = ValidateRoutes(config.Routes)
		if err != nil {
			return fmt.Errorf("Server %d: %s", n, err.Error())}

	}

	return nil
}

func ValidateErrorPage(path string) error {
	// Directory exists
	info, err := os.Stat(path)
		if err != nil {
			return err}
	// make sure it is a directory
	if !info.IsDir() {
			return errors.New("File is not a directory")
	}

	return nil
}

func ValidateRoutes(routes []Route) error {

	usedUrls := make([]string,0)

	for _, route := range routes {

		// has to start with a '/'
		if route.Url == "" || route.Url[0] != '/' {
			return errors.New("Route path must start with a /")}
		// can not contain a query string or a fragment
		if strings.Contains(route.Url, "?") || strings.Contains(route.Url, "#") {
			return errors.New("Route path can not contain a query string or a fragment")}
		// url has to be normalized
		if route.Url != path.Clean(route.Url) {
			return errors.New("Route path has to be normalized")}
		// handler must be supported by server
		if !slices.Contains(MyServerSettings.Supports, route.Handler) {
			return errors.New("Handler not supported by server")}
		// route prefix must not repeat
		if !slices.Contains(usedUrls, route.Url) {
		usedUrls = append(usedUrls, route.Url)
		} else {
			return errors.New("Route url is already in use")}
	}

	return nil
}

// 31. CGI-specific directives must not be required while CGI is not implemented.