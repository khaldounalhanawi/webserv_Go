package router

import (
	"fmt"
	"testing"
)

func	TestRouter(t *testing.T) {

	routesList := []Route{{"/", "GET"}, {"/abc", "GET"}, {"/abc/efg", "GET"}}
	
	request := Request{"GET", "/abc", "HTTP/1.1", map[string]string{"host":"BMW", "Content-Length":"5"}, []byte("hello")}

	router := Router(request, routesList)

	fmt.Println(router.url)
}