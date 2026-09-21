package types

type Request struct {
	Method	string
	Path	string
	Version	string
	Headers	map[string]string
	Body	[]byte
}

type Response struct {
	Status		int
	ContentType	string
	Body		[]byte
}

type Route struct {
	Url		string
	Handler	string
}

type Config struct {
	ListenPort		int
	Root			string
	Routes			[]Route
	Error_page		string
	Max_header_size	int
	Max_body_size	int
}

type ServerSettings struct {
	Maximum_header_limit	int
	Maximum_body_limit		int
	Default_body_size		int
	Default_header_size		int
	Supports				[]string
}