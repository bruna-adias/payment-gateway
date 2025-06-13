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
	"payment-gateway/cmd/infra/handler"
)

type MockProcessPaymentUseCase struct {
	mock.Mock
}

func (m *MockProcessPaymentUseCase) Execute(paymentID int64, paymentType string, details string) error {
	args := m.Called(paymentID, paymentType, details)
	return args.Error(0)
}

func setupProcessPaymentTestRouter(h *handler.ProcessPaymentHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/payments/:id/process", h.Execute)
	return r
}

func TestProcessPaymentHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(MockProcessPaymentUseCase)
	h := handler.NewProcessPaymentHandler(mockUC)
	r := setupProcessPaymentTestRouter(h)

	paymentID := int64(123)
	paymentType := "credit_card"
	details := "card ending in 4242"

	mockUC.On("Execute", paymentID, paymentType, details).Return(nil)

	reqBody := map[string]interface{}{
		"type":    paymentType,
		"details": details,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/payments/123/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "payment processed successfully", resp["message"])
	assert.Equal(t, float64(paymentID), resp["payment_id"])
	assert.Equal(t, paymentType, resp["type"])
	assert.Equal(t, details, resp["details"])
}

func TestProcessPaymentHandler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := handler.NewProcessPaymentHandler(nil)
	r := setupProcessPaymentTestRouter(h)

	req, _ := http.NewRequest(http.MethodPost, "/payments/invalid/process", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProcessPaymentHandler_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := handler.NewProcessPaymentHandler(nil)
	r := setupProcessPaymentTestRouter(h)

	body := bytes.NewBufferString("{invalid json}")
	req, _ := http.NewRequest(http.MethodPost, "/payments/123/process", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProcessPaymentHandler_UseCaseError_500(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(MockProcessPaymentUseCase)
	h := handler.NewProcessPaymentHandler(mockUC)
	r := setupProcessPaymentTestRouter(h)

	paymentID := int64(123)
	paymentType := "credit_card"
	details := "card ending in 4242"

	mockUC.On("Execute", paymentID, paymentType, details).Return(assert.AnError)

	reqBody := map[string]interface{}{
		"type":    paymentType,
		"details": details,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/payments/123/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestProcessPaymentHandler_UseCaseError_400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(MockProcessPaymentUseCase)
	h := handler.NewProcessPaymentHandler(mockUC)
	r := setupProcessPaymentTestRouter(h)

	paymentID := int64(123)
	paymentType := "credit_card"
	details := "card ending in 4242"

	mockUC.On("Execute", paymentID, paymentType, details).Return(exceptions.NewDomainError("error processing payment"))

	reqBody := map[string]interface{}{
		"type":    paymentType,
		"details": details,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/payments/123/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertExpectations(t)
}
