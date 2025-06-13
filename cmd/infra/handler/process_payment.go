package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	exceptions "payment-gateway/cmd/domain/err"
	"strconv"
)

type ProcessPaymentUseCase interface {
	Execute(paymentID int64, paymentType string, details string) error
}

type ProcessPaymentHandler struct {
	useCase ProcessPaymentUseCase
}

func NewProcessPaymentHandler(useCase ProcessPaymentUseCase) *ProcessPaymentHandler {
	return &ProcessPaymentHandler{
		useCase: useCase,
	}
}

func (h *ProcessPaymentHandler) Execute(ctx *gin.Context) {
	paymentID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid payment id"})
		return
	}

	var request struct {
		Type    string `json:"type"`
		Details string `json:"details"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	err = h.useCase.Execute(paymentID, request.Type, request.Details)
	if err != nil {
		var ex *exceptions.DomainError

		if errors.As(err, &ex) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": ex.Error()})
			return
		}

		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"message":    "payment processed successfully",
		"payment_id": paymentID,
		"type":       request.Type,
		"details":    request.Details,
	})
}
