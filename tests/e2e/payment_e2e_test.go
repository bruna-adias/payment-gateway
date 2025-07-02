//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var baseURL string

func init() {
	baseURL = os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
}

type PaymentRequest struct {
	OrderID     int64   `json:"order_id"`
	Amount      float64 `json:"amount"`
	PaymentType string  `json:"payment_type"`
}

type PaymentResponse struct {
	ID          int64  `json:"id"`
	OrderID     int64  `json:"order_id"`
	Status      string `json:"status"`
	PaymentType string `json:"type"`
}

type PaymentProcessedResponse struct {
	ID      int64  `json:"payment_id"`
	Type    string `json:"type"`
	Msg     string `json:"message"`
	Details string `json:"details"`
}

type ProcessPaymentRequest struct {
	Type    string `json:"type"`
	Details string `json:"details"`
}

type OrderResponse struct {
	ID      int64           `json:"id"`
	Status  string          `json:"status"`
	Amount  float64         `json:"amount"`
	Cashout CashoutResponse `json:"cashout"`
}

type CashoutResponse struct {
	CashedDebt    float64 `json:"cashed_debt"`
	RemainingDebt float64 `json:"remaining_debt"`
	Charges       float64 `json:"charges"`
	IsPaid        bool    `json:"is_paid"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestPaymentTotalPaidFlow(t *testing.T) {
	orderID := int64(1)
	amount := 120.5
	paymentType := "CreditCard"
	var paymentID int64

	t.Run("should create a payment", func(t *testing.T) {
		paymentReq := PaymentRequest{
			OrderID:     orderID,
			Amount:      amount,
			PaymentType: paymentType,
		}

		reqBody, err := json.Marshal(paymentReq)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/payments", baseURL)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var paymentResp PaymentResponse
		err = json.NewDecoder(resp.Body).Decode(&paymentResp)
		require.NoError(t, err)

		assert.NotZero(t, paymentResp.ID)
		assert.Equal(t, orderID, paymentResp.OrderID)
		assert.Equal(t, "pending", paymentResp.Status)
		assert.Equal(t, paymentType, paymentResp.PaymentType)

		paymentID = paymentResp.ID
	})

	t.Run("should process a payment", func(t *testing.T) {
		require.NotZero(t, paymentID, "paymentID should be set from create test")

		processReq := ProcessPaymentRequest{Type: "Success", Details: "approved details"}
		reqBody, err := json.Marshal(processReq)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/payments/%d/process", baseURL, paymentID)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var paymentResp PaymentProcessedResponse
		err = json.NewDecoder(resp.Body).Decode(&paymentResp)
		require.NoError(t, err)

		assert.Equal(t, "Success", paymentResp.Type)
		assert.Equal(t, "payment processed successfully", paymentResp.Msg)
		assert.Equal(t, "approved details", paymentResp.Details)
		assert.Equal(t, paymentID, paymentResp.ID)
	})

	t.Run("should get an order with payment details", func(t *testing.T) {
		require.NotZero(t, orderID, "orderID should be set")

		url := fmt.Sprintf("%s/orders/%d", baseURL, orderID)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var orderResp OrderResponse
		err = json.NewDecoder(resp.Body).Decode(&orderResp)
		require.NoError(t, err)

		assert.Equal(t, orderID, orderResp.ID)
		assert.Equal(t, 120.5, orderResp.Amount)
		assert.Equal(t, 120.5, orderResp.Cashout.CashedDebt)
		assert.Equal(t, 0.0, orderResp.Cashout.RemainingDebt)
		assert.Equal(t, 12.05, orderResp.Cashout.Charges)
		assert.Equal(t, true, orderResp.Cashout.IsPaid)
		assert.Equal(t, "paid", orderResp.Status)
	})
}

func TestPaymentNonTotalPaidFlow(t *testing.T) {
	orderID := int64(2)
	amount := 120.5
	paymentType := "CreditCard"
	var paymentID int64

	t.Run("should create a payment", func(t *testing.T) {
		paymentReq := PaymentRequest{
			OrderID:     orderID,
			Amount:      amount,
			PaymentType: paymentType,
		}

		reqBody, err := json.Marshal(paymentReq)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/payments", baseURL)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var paymentResp PaymentResponse
		err = json.NewDecoder(resp.Body).Decode(&paymentResp)
		require.NoError(t, err)

		assert.NotZero(t, paymentResp.ID)
		assert.Equal(t, orderID, paymentResp.OrderID)
		assert.Equal(t, "pending", paymentResp.Status)
		assert.Equal(t, paymentType, paymentResp.PaymentType)

		paymentID = paymentResp.ID
	})

	t.Run("should process a payment", func(t *testing.T) {
		require.NotZero(t, paymentID, "paymentID should be set from create test")

		processReq := ProcessPaymentRequest{Type: "Success", Details: "approved details"}
		reqBody, err := json.Marshal(processReq)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/payments/%d/process", baseURL, paymentID)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var paymentResp PaymentProcessedResponse
		err = json.NewDecoder(resp.Body).Decode(&paymentResp)
		require.NoError(t, err)

		assert.Equal(t, "Success", paymentResp.Type)
		assert.Equal(t, "payment processed successfully", paymentResp.Msg)
		assert.Equal(t, "approved details", paymentResp.Details)
		assert.Equal(t, paymentID, paymentResp.ID)
	})

	t.Run("should get an order with payment details", func(t *testing.T) {
		require.NotZero(t, orderID, "orderID should be set")

		url := fmt.Sprintf("%s/orders/%d", baseURL, orderID)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var orderResp OrderResponse
		err = json.NewDecoder(resp.Body).Decode(&orderResp)
		require.NoError(t, err)

		assert.Equal(t, orderID, orderResp.ID)
		assert.Equal(t, 250.0, orderResp.Amount)
		assert.Equal(t, 120.5, orderResp.Cashout.CashedDebt)
		assert.Equal(t, 129.5, orderResp.Cashout.RemainingDebt)
		assert.Equal(t, 12.05, orderResp.Cashout.Charges)
		assert.Equal(t, false, orderResp.Cashout.IsPaid)
		assert.Equal(t, "pending", orderResp.Status)
	})
}

func TestPaymentFlowError(t *testing.T) {
	orderID := int64(3)
	amount := 500.0
	paymentType := "CreditCard"

	t.Run("should create a payment", func(t *testing.T) {
		paymentReq := PaymentRequest{
			OrderID:     orderID,
			Amount:      amount,
			PaymentType: paymentType,
		}

		reqBody, err := json.Marshal(paymentReq)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/payments", baseURL)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errorResponse ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)

		assert.Equal(t, "Payment exceeds debt", errorResponse.Error)
	})
}
