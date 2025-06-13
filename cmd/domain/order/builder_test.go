package order_test

import (
	"payment-gateway/cmd/domain/order"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderBuilder(t *testing.T) {
	t.Run("should create new builder with empty order", func(t *testing.T) {
		b := order.NewOrderBuilder()
		assert.NotNil(t, b)
		assert.NotNil(t, b.Build())
	})
}

func TestBuilderMethods(t *testing.T) {
	now := time.Now()

	t.Run("should build order with all fields set", func(t *testing.T) {
		p := order.NewOrderBuilder().
			WithId(1).
			WithAmount(100.0).
			WithStatus("approved").
			WithCreatedAt(now).
			WithUpdatedAt(now).
			Build()

		assert.Equal(t, int64(1), p.Id())
		assert.Equal(t, 100.0, p.Amount())
		assert.Equal(t, "approved", p.Status())
		assert.Equal(t, now, p.CreatedAt())
		assert.Equal(t, now, p.UpdatedAt())
	})
}
