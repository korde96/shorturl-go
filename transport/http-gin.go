package transport

import (
	"fmt"
	"net/http"
	"shorturl-go/endpoint"

	"github.com/gin-gonic/gin"
)

//return http.Handler
func NewHTTPHandler(endpoints endpoint.Endpoints) http.Handler {
	r := gin.Default()

	r.POST("/add", func(c *gin.Context) {
		var req endpoint.AddURLRequest
		err := c.ShouldBind(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		resp, err := endpoints.AddURLEndpoint(c, req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, resp)

	})

	r.GET("/:surl", func(c *gin.Context) {
		surl := c.Param("surl")
		req := endpoint.GetURLRequest{SURL: surl}
		if resp, err := endpoints.GetURLEndpoint(c, req); err == nil {
			resp := resp.(endpoint.GetURLResponse)
			c.Redirect(http.StatusFound, fmt.Sprintf("//%s", resp.URL))
		} else {
			c.String(http.StatusNotFound, err.Error())
		}
	})

	// r.GET("/:surl", pkg.GinHandler(endpoints.GetURLEndpoint,
	// 	func(c *gin.Context) (request interface{}, err error) {
	// 		surl := c.Param("surl")
	// 		request = endpoint.GetURLRequest{SURL: surl}
	// 		return request, nil
	// 	},
	// 	func(c *gin.Context, response interface{}) {
	// 		resp := response.(endpoint.GetURLResponse)
	// 		c.Redirect(http.StatusFound, fmt.Sprintf("//%s", resp.URL))
	// 	},
	// 	func(c *gin.Context, err error) {
	// 		c.String(http.StatusNotFound, err.Error())
	// 	}))

	return r
}
