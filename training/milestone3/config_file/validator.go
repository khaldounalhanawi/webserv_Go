package LoadConfig

import (
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
)

func ValidateConfigs(configs []*Config) error {

	usedPorts := make ([]int,0)

	for n, config := range configs {

		// must have (listen, root, at least one route)
		if !config.declared["listen"] {return fmt.Errorf("Server %d: Missing Listen in Config file", n)}
		if !config.declared["root"] {return fmt.Errorf("Server %d: Missing root in Config file", n)}
		if len(config.routes) < 1 {return fmt.Errorf("Server %d: Must have atleast one Route", n)}

		// port: must be between 1 and 65535
		if config.ListenPort < 1 || config.ListenPort > 65535 {
			return fmt.Errorf("Server %d: Not a valid port number", n)}
		// port: must not use already used port
		if !slices.Contains(usedPorts, config.ListenPort) {
			usedPorts = append(usedPorts, config.ListenPort)
			} else {
				return fmt.Errorf("Server %d: Port %d is already in use", n, config.ListenPort)}

		// root: The root path must not be empty
		if config.root == "" {return fmt.Errorf("Server %d: Root path is empty", n)}
		// root: The root path must refer to an existing directory
		dir, err := os.Stat(config.root)
		if err != nil {return fmt.Errorf("Server %d: %s", n, err.Error())}
		if !dir.IsDir() {return fmt.Errorf("Server %d: Root is not a Directory", n)}

		if config.declared["error_page"] {
			err := ValidateErrorPage(config.error_page)
			if err != nil {
				return fmt.Errorf("Server %d: %s", n, err.Error())}
		}

		// Max header size must be between 0 and limit from server settings
		if config.declared["max_header_size"] {
			if config.max_header_size < 0 || config.max_header_size > MyServerSettings.maximum_header_limit {
				return fmt.Errorf("Server %d: Max header size must be between 0 and %d", n, MyServerSettings.maximum_header_limit)}
		} else {
			config.max_header_size = MyServerSettings.default_header_size }

		// Max body size must be between 0 and limit from server settings
		if config.declared["max_body_size"] {
			if config.max_body_size < 0 || config.max_body_size > MyServerSettings.maximum_body_limit {
				return fmt.Errorf("Server %d: Max body size must be between 0 and %d", n, MyServerSettings.maximum_body_limit)}
		} else {
			config.max_body_size = MyServerSettings.default_body_size }
		
		// Routes: validate routes
		err = ValidateRoutes(config.routes)
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
		if route.url == "" || route.url[0] != '/' {
			return errors.New("Route path must start with a /")}
		// can not contain a query string or a fragment
		if strings.Contains(route.url, "?") || strings.Contains(route.url, "#") {
			return errors.New("Route path can not contain a query string or a fragment")}
		// url has to be normalized
		if route.url != path.Clean(route.url) {
			return errors.New("Route path has to be normalized")}
		// handler must be supported by server
		if !slices.Contains(MyServerSettings.supports, route.handler) {
			return errors.New("Handler not supported by server")}
		// route prefix must not repeat
		if !slices.Contains(usedUrls, route.url) {
		usedUrls = append(usedUrls, route.url)
		} else {
			return errors.New("Route url is already in use")}
	}

	return nil
}

// func CheckOutsideRoot(path string, root string) error {

// 	rel, err := filepath.Rel(root, filepath.Join(root, path))
// 	if err != nil {
// 		return err}
// 	if rel == ".." || strings.HasPrefix(rel, ".." + string(filepath.Separator)) {
// 		return errors.New("Path goes outside root")}

// 	return nil
// }

// 32. All validation errors must identify the server and directive involved when possible.
// 1. The configuration must contain at least one server block.
// 3. The listen port must be an integer between 1 and 65535.
// 4. No two server blocks may use the same listen port.
// 5. Every server must declare "root" exactly once.
// 6. The root path must not be empty.
// 20. Every server should contain at least one route.
// 7. The root path must refer to an existing directory.
// 8. The root path must not refer to a regular file.
// 10. If "error_page" is declared, its path must not be empty.
// 11. If "error_page" is declared, its path must refer to an existing directory.
// 13. If "max_header_size" is omitted, use the default value of 8192 bytes.
// 14. If "max_header_size" is declared, it must be a positive integer.
// 15. "max_header_size" must not exceed the server's maximum allowed header limit.
// 17. If "max_body_size" is omitted, use the default value of 10485760 bytes.
// 18. If "max_body_size" is declared, it must be a positive integer.
// 19. "max_body_size" must not exceed the server's maximum allowed body limit.
// 21. Every route URL must begin with '/'.
// 22. Every route URL must be a valid URL path.
// 23. Route URLs must not contain a query string or fragment.
// 24. Route URLs must be normalized consistently.
// 29. A route prefix must not resolve outside the server's document root.
// 25. Every route handler must be supported by the server.
// 26. For the current project, supported handlers are "static" and, when implemented, "cgi".
// 27. Exact duplicate route prefixes must be rejected.


// 31. CGI-specific directives must not be required while CGI is not implemented.
