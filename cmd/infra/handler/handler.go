package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	Execute(ctx *gin.Context)
}
