package requestHandler

import (
	"errors"
	. "types"
)

func RequestHandler(request Request, config Config) (*Response, error) {

	// myroute := router.Router(request, config.Routes)

	// if myroute.Handler == "static" {
	// 	return StaticHandler(request, config)
	// } else if myroute.Handler == "cgi" {
	// 	return CgiHandler(request, config)
	// }

	return nil, errors.New("Not supported handler")
}