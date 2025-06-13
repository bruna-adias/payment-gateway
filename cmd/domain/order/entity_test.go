package order_test

import (
	"payment-gateway/cmd/domain/order"
	"payment-gateway/cmd/domain/payment"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEntitySetters(t *testing.T) {
	o := order.NewOrderBuilder().Build()

	t.Run("should set id correctly", func(t *testing.T) {
		o.SetId(456)

		assert.Equal(t, int64(456), o.Id())
	})

	t.Run("should set amount correctly", func(t *testing.T) {
		originalUpdatedAt := o.UpdatedAt()
		time.Sleep(100 * time.Millisecond)
		newAmount := 456.0
		o.SetAmount(newAmount)
		assert.Equal(t, newAmount, o.Amount())
		assert.True(t, o.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("should set status correctly", func(t *testing.T) {
		originalUpdatedAt := o.UpdatedAt()
		time.Sleep(100 * time.Millisecond)
		newStatus := "approved"
		o.SetStatus(newStatus)
		assert.Equal(t, newStatus, o.Status())
		assert.True(t, o.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("should set createdAt correctly", func(t *testing.T) {
		createdAt := time.Now()
		o.SetCreatedAt(createdAt)
		assert.Equal(t, createdAt, o.CreatedAt())
	})

	t.Run("should set updatedAt correctly", func(t *testing.T) {
		createdAt := time.Now()
		o.SetUpdateAt(createdAt)
		assert.Equal(t, createdAt, o.UpdatedAt())
	})
}

func TestEntityProcessPayment(t *testing.T) {
	o := order.NewOrderBuilder().WithStatus("pending").Build()
	p := payment.NewPaymentBuilder().WithAmount(123).Build()

	t.Run("should pay order", func(t *testing.T) {
		p.SetStatus("approved")
		err := o.ProcessPayment(123.0, *p)

		assert.NoError(t, err)
		assert.Equal(t, "paid", o.Status())
	})

	t.Run("should not pay order", func(t *testing.T) {
		p.SetStatus("approved")
		o.SetStatus("pending")

		err := o.ProcessPayment(150.0, *p)

		assert.NoError(t, err)
		assert.Equal(t, "pending", o.Status())
	})

	t.Run("should throw error to pay", func(t *testing.T) {
		p.SetStatus("approved")
		o.SetStatus("pending")

		err := o.ProcessPayment(100.0, *p)

		assert.Equal(t, "Payment exceeds debt", err.Error())
		assert.Equal(t, "pending", o.Status())
	})
}

func TestEntityPreValidation(t *testing.T) {
	o := order.NewOrderBuilder().WithStatus("pending").Build()

	t.Run("Should not throw error in validation", func(t *testing.T) {
		err := o.PreValidation(100.0, 90.0)
		assert.NoError(t, err)
	})

	t.Run("Should throw error when the payment is larger than the remaining debt", func(t *testing.T) {
		err := o.PreValidation(100.0, 110.0)
		assert.Equal(t, "Payment exceeds debt", err.Error())
	})
}
