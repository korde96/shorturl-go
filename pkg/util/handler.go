package util

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
)

func GinHandler(e endpoint.Endpoint,
	dec func(c *gin.Context) (request interface{}, err error),
	respond func(c *gin.Context, response interface{}),
	errorHandler func(c *gin.Context, err error),
) gin.HandlerFunc {

	return func(c *gin.Context) {
		request, err := dec(c)
		if err != nil {
			errorHandler(c, err)
			return
		}
		if response, err := e(c, request); err == nil {
			respond(c, response)
		} else {
			errorHandler(c, err)
		}
	}
}
