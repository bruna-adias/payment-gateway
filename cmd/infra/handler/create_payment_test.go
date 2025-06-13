package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	exceptions "payment-gateway/cmd/domain/err"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"payment-gateway/cmd/domain/payment"
	"payment-gateway/cmd/infra/handler"
)

type MockCreatePaymentUseCase struct {
	mock.Mock
}

func (m *MockCreatePaymentUseCase) Execute(orderId int64, amount float64, status string) (*payment.Entity, error) {
	args := m.Called(orderId, amount, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*payment.Entity), args.Error(1)
}

func setupTestRouter(h *handler.CreatePaymentHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/payments", h.Execute)
	return r
}

func TestCreatePaymentHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(MockCreatePaymentUseCase)
	h := handler.NewCreatePaymentHandler(mockUC)
	r := setupTestRouter(h)

	orderID := int64(123)
	amount := 100.50
	paymentType := "credit_card"
	expectedPayment := payment.NewPayment(orderID, amount, paymentType)
	expectedPayment.SetId(1)

	mockUC.On("Execute", orderID, amount, paymentType).Return(expectedPayment, nil)

	reqBody := map[string]interface{}{
		"order_id":     orderID,
		"payment_type": paymentType,
		"amount":       100.50,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp["id"])
	assert.Equal(t, float64(orderID), resp["order_id"])
	assert.Equal(t, "pending", resp["status"])
	assert.Equal(t, paymentType, resp["type"])
}

func TestCreatePaymentHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := handler.NewCreatePaymentHandler(nil)
	r := setupTestRouter(h)

	body := bytes.NewBufferString("{invalid json}")
	req, _ := http.NewRequest(http.MethodPost, "/payments", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePaymentHandler_UseCaseError_500(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(MockCreatePaymentUseCase)
	h := handler.NewCreatePaymentHandler(mockUC)
	r := setupTestRouter(h)

	orderID := int64(123)
	amount := 100.50
	paymentType := "credit_card"
	mockUC.On("Execute", orderID, amount, paymentType).Return(nil, assert.AnError)

	reqBody := map[string]interface{}{
		"order_id":     orderID,
		"payment_type": paymentType,
		"amount":       100.50,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreatePaymentHandler_UseCaseError_400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(MockCreatePaymentUseCase)
	h := handler.NewCreatePaymentHandler(mockUC)
	r := setupTestRouter(h)

	orderID := int64(123)
	amount := 100.50
	paymentType := "credit_card"
	mockUC.On("Execute", orderID, amount, paymentType).Return(nil, exceptions.NewDomainError("error creating payment"))

	reqBody := map[string]interface{}{
		"order_id":     orderID,
		"payment_type": paymentType,
		"amount":       100.50,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertExpectations(t)
}
