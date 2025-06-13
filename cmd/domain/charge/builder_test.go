package charge_test

import (
	"payment-gateway/cmd/domain/charge"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewChargeBuilder(t *testing.T) {
	t.Run("should create new builder with empty payment", func(t *testing.T) {
		b := charge.NewChargeBuilder()
		assert.NotNil(t, b)
		assert.NotNil(t, b.Build())
	})
}

func TestBuilderMethods(t *testing.T) {
	now := time.Now()

	t.Run("should build payment with all fields set", func(t *testing.T) {
		p := charge.NewChargeBuilder().
			WithId(1).
			WithAmount(100.0).
			WithCategory("financial_fee").
			WithPaymentId(123).
			WithCreatedAt(now).
			WithUpdatedAt(now).
			Build()

		assert.Equal(t, int64(1), p.Id())
		assert.Equal(t, "financial_fee", p.Category())
		assert.Equal(t, int64(123), p.PaymentId())
		assert.Equal(t, 100.0, p.Amount())
		assert.Equal(t, now, p.CreatedAt())
		assert.Equal(t, now, p.UpdatedAt())
	})
}
