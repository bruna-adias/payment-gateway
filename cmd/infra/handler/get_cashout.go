package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	exceptions "payment-gateway/cmd/domain/err"
	"payment-gateway/cmd/domain/order"
	"payment-gateway/cmd/usecases"
	"strconv"
)

type GetCashoutUseCase interface {
	Execute(orderId int64) (order.Entity, usecases.CashoutView, error)
}

type GetCashoutHandler struct {
	UseCase GetCashoutUseCase
}

func NewGetCashoutHandler(useCase GetCashoutUseCase) *GetCashoutHandler {
	return &GetCashoutHandler{
		UseCase: useCase,
	}
}

func (c *GetCashoutHandler) Execute(ctx *gin.Context) {
	orderID := ctx.Param("id")
	orderId64, err := strconv.ParseInt(orderID, 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	or, view, err := c.UseCase.Execute(orderId64)
	if err != nil {
		var ex *exceptions.DomainError

		if errors.As(err, &ex) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": ex.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":      or.Id(),
		"amount":  or.Amount(),
		"status":  or.Status(),
		"cashout": view,
	})
}
