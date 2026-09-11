package router

import (
	"fmt"
	"testing"
)

type test struct {
	answer	string
	input	string
}

func	TestRouter(t *testing.T) {

	tests := []test{
	// Exact matches
	{"/", "/"},
	{"/abc", "/abc"},
	{"/abc/efg", "/abc/efg"},
	{"/uploads", "/uploads"},

	// Trailing slash
	{"/abc", "/abc/"},
	{"/abc/efg", "/abc/efg/"},
	{"/uploads", "/uploads/"},

	// Nested paths
	{"/abc", "/abc/x"},
	{"/abc/efg", "/abc/efg/x"},
	{"/abc/efg", "/abc/efg/x/y"},
	{"/uploads", "/uploads/a/b/c.txt"},

	// Prefix boundary — must NOT match the longer word
	{"/", "/abcd"},
	{"/", "/abcde"},
	{"/", "/uploads2"},
	{"/", "/uploadsomething"},

	// Longest-match behavior
	{"/abc", "/abc/test"},
	{"/abc/efg", "/abc/efg/test"},
	{"/uploads", "/uploads/test.jpg"},

	// Similar-looking paths
	{"/", "/ab"},
	{"/", "/abcc"},
	{"/", "/upload"},
	}

	routesList := []Route{{"/", "GET"}, {"/abc/efg", "GET"}, {"/abc", "GET"}, {"/uploads", "GET"}}
	request := Request{"GET", "PLACEHOLDER", "HTTP/1.1", map[string]string{"host":"BMW", "Content-Length":"5"}, []byte("hello")}

	for i := range tests{
		request.Path = tests[i].input
		router := Router(request, routesList)
		if router.url == tests[i].answer {
			fmt.Printf("Passed test %-20s --> %-20s\n", tests[i].input, tests[i].answer)
		} else {
			fmt.Printf("Failed test %-20s --> %-20s\n", tests[i].input, tests[i].answer)
			t.Errorf("Failed test %-20s --> %-20s\n", tests[i].input, tests[i].answer)
		}
	}
}