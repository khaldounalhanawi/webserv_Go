package main

import (
	"fmt"
	"strings"
)

// func main() {
//   incomeUrl := "/abc/efg/falafel.pic"
//   serverSettings := []string{"/", "/abc","/abc/efg","/uploads"}
//   matchIndex := -1

//   incomeUrl = strings.TrimRight(incomeUrl, "/")

//   array := strings.Split(incomeUrl, "/")
//   for n, item := range array {
//     fmt.Printf("item%d: %s\n", n,item)
//   }

//   compareLine := "/"
//   for i := 0; i < len(array); i++ {
//     compareLine += array[i]
//     if slices.Contains(serverSettings, compareLine) {
//       matchIndex = i
//     }
//     if matchIndex != 0 {compareLine += "/"}
//   }

//   fmt.Println("Final route is ",serverSettings[matchIndex])

// }
type Route struct {
	Path    string
	Handler string
}

type Request struct {
	Method	string
	Path	string
	Version	string
	Headers	map[string]string
	Body	[]byte
}

func Router(request Request, routes []Route) *Route {
	var bestMatch *Route

	for i := range routes {
		route := &routes[i]

		if !matches(route.Path, request.Path) {
			continue
		}

		if bestMatch == nil || len(route.Path) > len(bestMatch.Path) {
			bestMatch = route
		}
	}

	return bestMatch
}

func matches(routePath string, requestPath string) bool {
  
  if routePath == "/" {
      return true
    }

  if requestPath == routePath {
		return true
	}

  result := strings.HasPrefix(requestPath, routePath+"/")
	return result
}

func main () {
    routes := []Route{
    {Path: "/", Handler: "static"},
    {Path: "/abc/efg", Handler: "cgi"},
    {Path: "/abc", Handler: "static"},
    }

    request := Request{
    Method: "GET",
    Path:   "/abc/efg/file.txt",
    // Path:   "/a",
    }

    route := Router(request, routes)
    if route != nil {
      fmt.Println((*route).Path)
    }
    

}