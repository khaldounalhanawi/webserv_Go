package gethandler

type Request struct {
	Method	string
	Path	string
	Version	string
	Headers	map[string]string
	Body	[]byte
}

type Route struct {
	url		string
	handler	string
}

type Config struct {
	ListenPort		int
	root			string
	routes			[]Route
	error_page		string
	max_header_size	int
	max_body_size	int

	declared		map[string]bool
}

type ServerSettings struct {
	maximum_header_limit	int
	maximum_body_limit		int
	default_body_size		int
	default_header_size		int
	supports				[]string
}