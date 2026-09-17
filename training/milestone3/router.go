package router

import (
	"strings"
)

func matches(route string, request string) bool {
	if route == "/" || route == request {
		return true
	}

	return strings.HasPrefix(request, route + "/")
}

func Router(request Request, routesList []Route) *Route {

	var final *Route

	for i := range routesList {
		item := &routesList[i]
		if matches(item.url, request.Path) {
			if final == nil || len(item.url) > len(final.url) {
				final = item
			}
		}
	}

	return final
}