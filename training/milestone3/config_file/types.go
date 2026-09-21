package loadConfig

type Route struct {
	url		string
	handler	string
}

type ParsedConfig struct {
	config		*Config
	declared	map[string]bool
}

type Config struct {
	ListenPort		int
	root			string
	routes			[]Route
	error_page		string
	max_header_size	int
	max_body_size	int
}

type ServerSettings struct {
	maximum_header_limit	int
	maximum_body_limit		int
	default_body_size		int
	default_header_size		int
	supports				[]string
}

var MyServerSettings = ServerSettings{
	maximum_header_limit:	8192,
	maximum_body_limit:		1048576,
	default_body_size:		1024,
	default_header_size:	512,
	supports:				[]string{"static", "cgi"},
}