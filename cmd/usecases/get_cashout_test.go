package usecases_test

import (
	"payment-gateway/cmd/domain/charge"
	"payment-gateway/cmd/domain/order"
	helpers_test "payment-gateway/cmd/testhelpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"payment-gateway/cmd/domain/payment"
	"payment-gateway/cmd/usecases"
)

func TestGetCashout_Execute(t *testing.T) {
	mockPaymentDao := new(helpers_test.MockPaymentDao)
	mockOrderDao := new(helpers_test.MockOrderDao)
	mockChargeDao := new(helpers_test.MockChargeDao)

	getCashoutUseCase := usecases.NewGetCashout(mockPaymentDao, mockOrderDao, mockChargeDao)

	t.Run("should get cashout", func(t *testing.T) {
		expectedOrder := order.NewOrderBuilder().WithId(1).WithAmount(100).Build()
		mockOrderDao.On("FindById", int64(1)).Return(expectedOrder, nil).Once()
		mockPaymentDao.On("FindByOrderId", int64(1)).Return([]payment.Entity{
			*payment.NewPaymentBuilder().WithStatus("approved").WithAmount(10).Build(),
			*payment.NewPaymentBuilder().WithStatus("approved").WithAmount(10).Build(),
		}, nil).Once()
		mockChargeDao.On("FindByOrderId", int64(1)).Return([]charge.Entity{
			*charge.NewChargeBuilder().WithAmount(5).Build(),
			*charge.NewChargeBuilder().WithAmount(5).Build(),
		}, nil).Once()

		or, view, err := getCashoutUseCase.Execute(1)

		assert.Equal(t, *expectedOrder, or)
		assert.Equal(t, usecases.CashoutView{
			OrderId:       1,
			CashedDebt:    20,
			RemainingDebt: 80,
			Charges:       10,
			IsPaid:        false,
		}, view)
		assert.Nil(t, err)
	})

	t.Run("should throw error when order not found", func(t *testing.T) {
		expectedOrder := order.NewOrderBuilder().WithId(1).WithAmount(100).Build()
		mockOrderDao.On("FindById", int64(1)).Return(expectedOrder, assert.AnError).Once()

		or, view, err := getCashoutUseCase.Execute(1)

		assert.Equal(t, order.Entity{}, or)
		assert.Equal(t, usecases.CashoutView{}, view)
		assert.Error(t, err)
	})

	t.Run("should thrown an error when payments find thrown error", func(t *testing.T) {
		expectedOrder := order.NewOrderBuilder().WithId(1).WithAmount(100).Build()
		mockOrderDao.On("FindById", int64(1)).Return(expectedOrder, nil).Once()
		mockPaymentDao.On("FindByOrderId", int64(1)).Return([]payment.Entity{}, assert.AnError).Once()

		or, view, err := getCashoutUseCase.Execute(1)

		assert.Equal(t, order.Entity{}, or)
		assert.Equal(t, usecases.CashoutView{}, view)
		assert.Error(t, err)
	})

	t.Run("should thrown an error when charges find thrown error", func(t *testing.T) {
		expectedOrder := order.NewOrderBuilder().WithId(1).WithAmount(100).Build()
		mockOrderDao.On("FindById", int64(1)).Return(expectedOrder, nil).Once()
		mockPaymentDao.On("FindByOrderId", int64(1)).Return([]payment.Entity{
			*payment.NewPaymentBuilder().WithStatus("approved").WithAmount(10).Build(),
			*payment.NewPaymentBuilder().WithStatus("approved").WithAmount(10).Build(),
		}, nil).Once()
		mockChargeDao.On("FindByOrderId", int64(1)).Return([]charge.Entity{}, assert.AnError).Once()

		or, view, err := getCashoutUseCase.Execute(1)

		assert.Equal(t, order.Entity{}, or)
		assert.Equal(t, usecases.CashoutView{}, view)
		assert.Error(t, err)
	})
}
