package charge_test

import (
	"github.com/stretchr/testify/assert"
	"payment-gateway/cmd/domain/charge"
	"payment-gateway/cmd/domain/payment"
	"testing"
	"time"
)

func TestNewCharge(t *testing.T) {
	paymentEntity := payment.NewPaymentBuilder().WithId(1).WithOrderId(123).WithStatus("approved").WithType("CreditCard").WithAmount(10.0).WithDetails("details").Build()

	t.Run("should create charge with correct values for financial_fee", func(t *testing.T) {
		chargeEntity, ok := charge.NewCharge(*paymentEntity)

		assert.True(t, ok)
		assert.Equal(t, 1.0, chargeEntity.Amount())
		assert.Equal(t, "financial_fee", chargeEntity.Category())
		assert.Equal(t, int64(1), chargeEntity.PaymentId())
		assert.NotZero(t, chargeEntity.CreatedAt())
		assert.NotZero(t, chargeEntity.UpdatedAt())
	})

	t.Run("should create charge with correct values for process_fee", func(t *testing.T) {
		paymentEntity.SetType("CashSlip")
		chargeEntity, ok := charge.NewCharge(*paymentEntity)

		assert.True(t, ok)
		assert.Equal(t, 2.0, chargeEntity.Amount())
		assert.Equal(t, "process_fee", chargeEntity.Category())
	})

	t.Run("should create charge with correct values for free type", func(t *testing.T) {
		paymentEntity.SetType("invalid_type")
		chargeEntity, ok := charge.NewCharge(*paymentEntity)

		assert.True(t, ok)
		assert.Equal(t, 0.0, chargeEntity.Amount())
		assert.Equal(t, "free", chargeEntity.Category())
	})

	t.Run("should create charge with correct values for cash", func(t *testing.T) {
		paymentEntity.SetType("Cash")
		chargeEntity, ok := charge.NewCharge(*paymentEntity)

		assert.True(t, ok)
		assert.Equal(t, 0.0, chargeEntity.Amount())
		assert.Equal(t, "free", chargeEntity.Category())
	})

	t.Run("should not create if payment status is invalid", func(t *testing.T) {
		paymentEntity.SetStatus("pending")
		chargeEntity, ok := charge.NewCharge(*paymentEntity)

		assert.False(t, ok)
		assert.Nil(t, chargeEntity)
	})
}

func TestEntitySetters(t *testing.T) {
	paymentEntity := payment.NewPaymentBuilder().WithId(1).WithOrderId(123).WithStatus("approved").WithType("CreditCard").WithAmount(10.0).WithDetails("details").Build()
	chargeEntity, _ := charge.NewCharge(*paymentEntity)

	t.Run("should set amount correctly", func(t *testing.T) {
		originalUpdatedAt := chargeEntity.UpdatedAt()
		time.Sleep(100 * time.Millisecond)
		newAmount := 15.0
		chargeEntity.SetAmount(newAmount)
		assert.Equal(t, newAmount, chargeEntity.Amount())
		assert.True(t, chargeEntity.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("should set category correctly", func(t *testing.T) {
		originalUpdatedAt := chargeEntity.UpdatedAt()
		time.Sleep(100 * time.Millisecond)
		newCategory := "new_category"
		chargeEntity.SetCategory(newCategory)
		assert.Equal(t, newCategory, chargeEntity.Category())
		assert.True(t, chargeEntity.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("should set payment ID correctly", func(t *testing.T) {
		originalUpdatedAt := chargeEntity.UpdatedAt()
		time.Sleep(100 * time.Millisecond)
		newPaymentId := int64(2)
		chargeEntity.SetPaymentId(newPaymentId)
		assert.Equal(t, newPaymentId, chargeEntity.PaymentId())
		assert.True(t, chargeEntity.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("should set ID correctly", func(t *testing.T) {
		id := int64(1)
		chargeEntity.SetId(id)
		assert.Equal(t, id, chargeEntity.Id())
	})

	t.Run("should set createdAt correctly", func(t *testing.T) {
		createdAt := time.Now()
		chargeEntity.SetCreatedAt(createdAt)
		assert.Equal(t, createdAt, chargeEntity.CreatedAt())
	})

	t.Run("should set updatedAt correctly", func(t *testing.T) {
		createdAt := time.Now()
		chargeEntity.SetUpdatedAt(createdAt)
		assert.Equal(t, createdAt, chargeEntity.UpdatedAt())
	})
}
