package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	exceptions "payment-gateway/cmd/domain/err"
	"payment-gateway/cmd/domain/payment"
)

type UseCase interface {
	Execute(orderId int64, amount float64, status string) (*payment.Entity, error)
}

type CreatePaymentHandler struct {
	UseCase UseCase
}

func NewCreatePaymentHandler(useCase UseCase) *CreatePaymentHandler {
	return &CreatePaymentHandler{
		UseCase: useCase,
	}
}

func (c *CreatePaymentHandler) Execute(ctx *gin.Context) {
	var request struct {
		OrderID     int64   `json:"order_id"`
		Amount      float64 `json:"amount"`
		PaymentType string  `json:"payment_type"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pay, err := c.UseCase.Execute(request.OrderID, request.Amount, request.PaymentType)
	if err != nil {
		var ex *exceptions.DomainError

		if errors.As(err, &ex) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": ex.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":       pay.Id(),
		"order_id": pay.OrderID(),
		"status":   pay.Status(),
		"type":     pay.Type(),
	})
}
