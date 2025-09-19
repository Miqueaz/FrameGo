package router

import "github.com/gin-gonic/gin"

type GroupRouter struct {
	Router *gin.RouterGroup
}

type AppRouter struct {
	*gin.Engine
}

type RouterInterface interface {
	Use(mw gin.HandlerFunc)
	GET(path string, fn any)
	POST(path string, fn any)
	PUT(path string, fn any)
	DELETE(path string, fn any)
	Execute(addr string) error
}
