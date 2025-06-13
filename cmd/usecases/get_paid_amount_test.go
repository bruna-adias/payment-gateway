package usecases_test

import (
	"github.com/stretchr/testify/assert"
	"payment-gateway/cmd/domain/payment"
	"payment-gateway/cmd/testhelpers"
	"payment-gateway/cmd/usecases"
	"testing"
)

func TestGetPaidAmount(t *testing.T) {
	mockPaymentDao := new(testhelpers.MockPaymentDao)

	t.Run("should get paid amount when all payment is paid", func(t *testing.T) {
		mockPaymentDao.On("FindByOrderId", int64(1)).Return([]payment.Entity{
			*payment.NewPaymentBuilder().WithAmount(100).WithStatus("approved").Build(),
			*payment.NewPaymentBuilder().WithAmount(100).WithStatus("approved").Build(),
		}, nil).Once()

		amount, err := usecases.GetPaidAmount(mockPaymentDao, 1)

		assert.NoError(t, err)
		assert.Equal(t, 200.0, amount)
	})

	t.Run("should get paid amount when some payment is paid", func(t *testing.T) {
		mockPaymentDao.On("FindByOrderId", int64(1)).Return([]payment.Entity{
			*payment.NewPaymentBuilder().WithAmount(100).WithStatus("approved").Build(),
			*payment.NewPaymentBuilder().WithAmount(100).WithStatus("reproved").Build(),
		}, nil).Once()

		amount, err := usecases.GetPaidAmount(mockPaymentDao, 1)

		assert.NoError(t, err)
		assert.Equal(t, 100.0, amount)
	})

	t.Run("should get paid amount when none payment is paid", func(t *testing.T) {
		mockPaymentDao.On("FindByOrderId", int64(1)).Return([]payment.Entity{
			*payment.NewPaymentBuilder().WithAmount(100).WithStatus("reproved").Build(),
			*payment.NewPaymentBuilder().WithAmount(100).WithStatus("reproved").Build(),
		}, nil).Once()

		amount, err := usecases.GetPaidAmount(mockPaymentDao, 1)

		assert.NoError(t, err)
		assert.Equal(t, 0.0, amount)
	})

	t.Run("should throw error when order not found", func(t *testing.T) {
		mockPaymentDao.On("FindByOrderId", int64(1)).Return(nil, assert.AnError).Once()

		amount, err := usecases.GetPaidAmount(mockPaymentDao, 1)

		assert.Error(t, err)
		assert.Equal(t, 0.0, amount)
	})
}
