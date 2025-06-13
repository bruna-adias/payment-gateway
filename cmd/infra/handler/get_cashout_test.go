package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payment-gateway/cmd/domain/order"
	"payment-gateway/cmd/usecases"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"payment-gateway/cmd/infra/handler"
)

type MockGetCheckoutUseCase struct {
	mock.Mock
}

func (m *MockGetCheckoutUseCase) Execute(orderId int64) (order.Entity, usecases.CashoutView, error) {

	args := m.Called(orderId)
	if args.Get(0) == nil {
		return order.Entity{}, usecases.CashoutView{}, args.Error(1)
	}
	return args.Get(0).(order.Entity), args.Get(1).(usecases.CashoutView), args.Error(2)
}

func setupGetCashoutTestRouter(h *handler.GetCashoutHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/orders/:id", h.Execute)
	return r
}

func TestGetCashoutHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(MockGetCheckoutUseCase)
	h := handler.NewGetCashoutHandler(mockUC)
	r := setupGetCashoutTestRouter(h)

	orderID := int64(123)
	orderExpected := *order.NewOrderBuilder().WithId(1).WithStatus("approved").WithAmount(100).Build()
	cashoutExpected := usecases.CashoutView{
		OrderId:       orderID,
		CashedDebt:    0,
		RemainingDebt: 100,
		Charges:       0,
		IsPaid:        false,
	}

	mockUC.On("Execute", orderID).Return(orderExpected, cashoutExpected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/orders/123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	cashoutView := resp["cashout"].(map[string]interface{})

	assert.Equal(t, float64(100), resp["amount"])
	assert.Equal(t, float64(1), resp["id"])
	assert.Equal(t, "approved", resp["status"])
	assert.Equal(t, cashoutExpected.IsPaid, cashoutView["is_paid"])
	assert.Equal(t, cashoutExpected.Charges, cashoutView["charges"])
	assert.Equal(t, cashoutExpected.RemainingDebt, cashoutView["remaining_debt"])
	assert.Equal(t, cashoutExpected.CashedDebt, cashoutView["cashed_debt"])
}
