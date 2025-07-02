package conf

import (
	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine, run *Runtime) {
	engine.POST("/payments", run.CreatePaymentHandler.Execute)
	engine.POST("/payments/:id/process", run.ProcessPaymentHandler.Execute)
	engine.GET("/orders/:id", run.GetCashoutHandler.Execute)
	engine.GET("/health", HealthHandler())
}

func HealthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	}
}
