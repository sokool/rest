package rest

import (
	"github.com/gin-gonic/gin"
	"go-rl/view"
	"reflect"
)

func Renderer(c *gin.Context) {
	c.Next()
	var p any
	var r view.Renderer
	var ok bool
	if p, ok = c.Get(responseBody); !ok {
		return
	}
	if p == nil || p == "" || reflect.ValueOf(p).IsZero() {
		return
	}
	if r, ok = p.(view.Renderer); !ok {
		return
	}

	b, err := view.JSON(r, c.Request.Header.Values("render")...)
	if err != nil {
		c.Set(responseError, err)
		return
	}

	c.Set(responseBody, b.JSON())
}
