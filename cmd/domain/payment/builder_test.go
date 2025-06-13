package payment_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"payment-gateway/cmd/domain/payment"
)

func TestNewPaymentBuilder(t *testing.T) {
	t.Run("should create new builder with empty payment", func(t *testing.T) {
		b := payment.NewPaymentBuilder()
		assert.NotNil(t, b)
		assert.NotNil(t, b.Build())
	})
}

func TestBuilderMethods(t *testing.T) {
	now := time.Now()

	t.Run("should build payment with all fields set", func(t *testing.T) {
		p := payment.NewPaymentBuilder().
			WithId(1).
			WithOrderId(123).
			WithAmount(100.0).
			WithType("credit_card").
			WithStatus("approved").
			WithDetails("test details").
			WithCreatedAt(now).
			Build()

		assert.Equal(t, int64(1), p.Id())
		assert.Equal(t, int64(123), p.OrderID())
		assert.Equal(t, 100.0, p.Amount())
		assert.Equal(t, "credit_card", p.Type())
		assert.Equal(t, "approved", p.Status())
		assert.Equal(t, "test details", p.Details())
		assert.Equal(t, now, p.CreatedAt())
	})

	t.Run("should chain builder methods correctly", func(t *testing.T) {
		b := payment.NewPaymentBuilder().
			WithOrderId(123).
			WithAmount(100.0).
			WithType("credit_card")

		p := b.Build()

		assert.Equal(t, int64(123), p.OrderID())
		assert.Equal(t, 100.0, p.Amount())
		assert.Equal(t, "credit_card", p.Type())
		assert.Equal(t, "pending", p.Status()) // default value
	})

	t.Run("should handle zero values correctly", func(t *testing.T) {
		p := payment.NewPaymentBuilder().Build()

		assert.Zero(t, p.Id())
		assert.Zero(t, p.OrderID())
		assert.Zero(t, p.Amount())
		assert.Empty(t, p.Type())
		assert.Equal(t, "pending", p.Status()) // default value
		assert.Empty(t, p.Details())
		assert.NotZero(t, p.CreatedAt()) // should be set automatically
	})
}
