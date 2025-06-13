package usecases_test

import (
	exceptions "payment-gateway/cmd/domain/err"
	"payment-gateway/cmd/domain/order"
	helpers_test "payment-gateway/cmd/testhelpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"payment-gateway/cmd/domain/payment"
	"payment-gateway/cmd/usecases"
)

func TestCreatePayment_Execute(t *testing.T) {
	orderID := int64(123)
	amount := 100.5
	paymentType := "credit_card"
	expectedPayment := payment.NewPayment(orderID, amount, paymentType)
	expectedOrder := order.NewOrderBuilder().WithId(orderID).WithAmount(100.5).Build()
	expectedPayment.SetId(1)

	t.Run("should create payment successfully with no existing payments", func(t *testing.T) {
		mockPaymentDao := new(helpers_test.MockPaymentDao)
		mockOrderDao := new(helpers_test.MockOrderDao)
		var existingPayments []payment.Entity

		mockPaymentDao.On("Insert", mock.Anything).Return(expectedPayment, nil)
		mockPaymentDao.On("FindByOrderId", mock.Anything, mock.Anything).Return(existingPayments, nil)
		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, nil)

		useCase := usecases.NewCreatePayment(mockPaymentDao, mockOrderDao)
		result, err := useCase.Execute(orderID, amount, paymentType)

		assert.NoError(t, err)
		assert.Equal(t, expectedPayment, result)
		mockPaymentDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})

	t.Run("should not create payment when order amount is exceeded", func(t *testing.T) {
		mockPaymentDao := new(helpers_test.MockPaymentDao)
		mockOrderDao := new(helpers_test.MockOrderDao)
		var existingPayments []payment.Entity

		expectedOrder := order.NewOrderBuilder().WithId(orderID).WithAmount(10.5).Build()
		expectedErr := exceptions.NewDomainError("Payment exceeds debt")

		mockPaymentDao.On("FindByOrderId", mock.Anything, mock.Anything).Return(existingPayments, nil)
		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, nil)

		useCase := usecases.NewCreatePayment(mockPaymentDao, mockOrderDao)
		result, err := useCase.Execute(orderID, amount, paymentType)

		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
		mockPaymentDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})

	t.Run("should not create payment when order left debt is exceeded", func(t *testing.T) {
		mockPaymentDao := new(helpers_test.MockPaymentDao)
		mockOrderDao := new(helpers_test.MockOrderDao)
		existPay := payment.NewPaymentBuilder().WithOrderId(orderID).WithAmount(100.5).WithId(1).WithStatus("approved").Build()
		existingPayments := []payment.Entity{*existPay}

		expectedErr := exceptions.NewDomainError("Payment exceeds debt")

		mockPaymentDao.On("FindByOrderId", mock.Anything, mock.Anything).Return(existingPayments, nil)
		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, nil)

		useCase := usecases.NewCreatePayment(mockPaymentDao, mockOrderDao)
		result, err := useCase.Execute(orderID, amount, paymentType)

		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
		mockPaymentDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})

	t.Run("should return error when paymentDao fails", func(t *testing.T) {
		mockPaymentDao := new(helpers_test.MockPaymentDao)
		mockOrderDao := new(helpers_test.MockOrderDao)
		var existingPayments []payment.Entity

		mockPaymentDao.On("Insert", mock.Anything).Return(nil, assert.AnError)
		mockPaymentDao.On("FindByOrderId", mock.Anything, mock.Anything).Return(existingPayments, nil)
		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, nil)

		useCase := usecases.NewCreatePayment(mockPaymentDao, mockOrderDao)
		result, err := useCase.Execute(orderID, amount, paymentType)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockPaymentDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})

	t.Run("should return error when paymentDao fails retrieving existing payments", func(t *testing.T) {
		mockPaymentDao := new(helpers_test.MockPaymentDao)
		mockOrderDao := new(helpers_test.MockOrderDao)

		mockPaymentDao.On("FindByOrderId", mock.Anything, mock.Anything).Return(nil, assert.AnError)
		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, nil)

		useCase := usecases.NewCreatePayment(mockPaymentDao, mockOrderDao)
		result, err := useCase.Execute(orderID, amount, paymentType)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockPaymentDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})
	t.Run("should return error when orderDao fails", func(t *testing.T) {
		mockPaymentDao := new(helpers_test.MockPaymentDao)
		mockOrderDao := new(helpers_test.MockOrderDao)

		mockOrderDao.On("FindById", mock.Anything).Return(nil, assert.AnError)

		useCase := usecases.NewCreatePayment(mockPaymentDao, mockOrderDao)
		result, err := useCase.Execute(orderID, amount, paymentType)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockPaymentDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})
}
