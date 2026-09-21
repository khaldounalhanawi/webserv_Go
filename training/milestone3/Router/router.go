package router

import (
	"strings"
	. "types"
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
		if matches(item.Url, request.Path) {
			if final == nil || len(item.Url) > len(final.Url) {
				final = item
			}
		}
	}

	return final
}