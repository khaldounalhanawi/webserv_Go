package router

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