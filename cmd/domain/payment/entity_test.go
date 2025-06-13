package payment_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"payment-gateway/cmd/domain/payment"
)

func TestNewPayment(t *testing.T) {
	t.Run("should create new payment with correct initial values", func(t *testing.T) {
		orderID := int64(123)
		amount := float64(123)
		paymentType := "credit_card"
		now := time.Now()

		p := payment.NewPayment(orderID, amount, paymentType)

		assert.Equal(t, orderID, p.OrderID())
		assert.Equal(t, paymentType, p.Type())
		assert.Equal(t, amount, p.Amount())
		assert.Equal(t, "pending", p.Status())
		assert.Equal(t, "", p.Details())
		assert.WithinDuration(t, now, p.CreatedAt(), time.Second)
		assert.WithinDuration(t, now, p.UpdatedAt(), time.Second)
	})
}

func TestProcess(t *testing.T) {
	t.Run("should approve payment when process type is Success", func(t *testing.T) {
		p := payment.NewPayment(123, 123.0, "credit_card")
		initialUpdatedAt := p.UpdatedAt()

		time.Sleep(time.Millisecond)
		p.Process("Success", "approved details")

		assert.Equal(t, "approved", p.Status())
		assert.Equal(t, "approved details", p.Details())
		assert.True(t, p.UpdatedAt().After(initialUpdatedAt))
	})

	t.Run("should reprove payment for any other process type", func(t *testing.T) {
		p := payment.NewPayment(123, 123.0, "credit_card")
		initialUpdatedAt := p.UpdatedAt()

		time.Sleep(time.Millisecond)
		p.Process("Failure", "reproved details")

		assert.Equal(t, "reproved", p.Status())
		assert.Equal(t, "reproved details", p.Details())
		assert.True(t, p.UpdatedAt().After(initialUpdatedAt))
	})
}

func TestEntitySetters(t *testing.T) {
	p := payment.NewPayment(123, 123.0, "credit_card")

	t.Run("should set payment id correctly", func(t *testing.T) {
		p.SetId(456)

		assert.Equal(t, int64(456), p.Id())
	})

	t.Run("should set order ID correctly", func(t *testing.T) {
		originalUpdatedAt := p.UpdatedAt()
		time.Sleep(100 * time.Millisecond)
		newOrderID := int64(456)
		p.SetOrderID(newOrderID)
		assert.Equal(t, newOrderID, p.OrderID())
		assert.True(t, p.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("should set amount correctly", func(t *testing.T) {
		originalUpdatedAt := p.UpdatedAt()
		time.Sleep(100 * time.Millisecond)
		newAmount := 456.0
		p.SetAmount(newAmount)
		assert.Equal(t, newAmount, p.Amount())
		assert.True(t, p.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("should set type correctly", func(t *testing.T) {
		originalUpdatedAt := p.UpdatedAt()
		time.Sleep(100 * time.Millisecond)
		newType := "cash_slip"
		p.SetType(newType)
		assert.Equal(t, newType, p.Type())
		assert.True(t, p.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("should set status correctly", func(t *testing.T) {
		originalUpdatedAt := p.UpdatedAt()
		time.Sleep(100 * time.Millisecond)
		newStatus := "approved"
		p.SetStatus(newStatus)
		assert.Equal(t, newStatus, p.Status())
		assert.True(t, p.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("should set details correctly", func(t *testing.T) {
		originalUpdatedAt := p.UpdatedAt()
		time.Sleep(100 * time.Millisecond)
		newDetails := "new details"
		p.SetDetails(newDetails)
		assert.Equal(t, newDetails, p.Details())
		assert.True(t, p.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("should set createdAt correctly", func(t *testing.T) {
		createdAt := time.Now()
		p.SetCreatedAt(createdAt)
		assert.Equal(t, createdAt, p.CreatedAt())
	})
}

func TestInvalid(t *testing.T) {
	p := payment.NewPayment(123, 123.0, "credit_card")
	t.Run("should be valid when status is approved", func(t *testing.T) {
		p.SetStatus("approved")
		assert.True(t, p.IsValid())
	})

	t.Run("should be invalid when status is reproved", func(t *testing.T) {
		p.SetStatus("reproved")
		assert.False(t, p.IsValid())
	})

	t.Run("should be invalid when status is pending", func(t *testing.T) {
		p.SetStatus("pending")
		assert.False(t, p.IsValid())
	})
}
