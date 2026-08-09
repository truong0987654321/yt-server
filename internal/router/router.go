package router

import "github.com/gin-gonic/gin"

type Dependencies struct{}

func Setup(deps Dependencies) *gin.Engine {
	r := gin.Default()
	return r
}
