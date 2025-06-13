package main

import (
	"github.com/gin-gonic/gin"
	"payment-gateway/cmd/infra"
	"payment-gateway/cmd/infra/conf"
)

func main() {
	r := gin.Default()

	c := infra.NewConfiguration()
	run := conf.NewRuntime(c)
	conf.Routes(r, run)

	r.Run(":8080")
}
