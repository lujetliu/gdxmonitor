package routers

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	v1 := r.Group("/v1")
	{
		v1.GET("/gdx/blocks", nil)
	}

	return r
}
