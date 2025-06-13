package usecases_test

import (
	"payment-gateway/cmd/domain/order"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"payment-gateway/cmd/domain/charge"
	"payment-gateway/cmd/domain/payment"
	"payment-gateway/cmd/testhelpers"
	"payment-gateway/cmd/usecases"
)

func TestProcessPayment_Execute(t *testing.T) {
	paymentID := int64(123)
	orderID := int64(456)
	processType := "Success"
	details := "payment processed"

	expectedOrder := order.NewOrderBuilder().WithId(orderID).WithAmount(100.5).Build()
	existingPayment := payment.NewPayment(orderID, 100.5, "credit_card")
	existingPayment.SetId(paymentID)

	t.Run("should process payment successfully", func(t *testing.T) {
		mockPaymentDao := new(testhelpers.MockPaymentDao)
		mockChargeDao := new(testhelpers.MockChargeDao)
		mockOrderDao := new(testhelpers.MockOrderDao)
		var existingPayments []payment.Entity

		mockPaymentDao.On("FindById", paymentID).Return(existingPayment, nil)
		mockPaymentDao.On("FindByOrderId", mock.Anything, mock.Anything).Return(existingPayments, nil)
		mockPaymentDao.On("Update", mock.Anything).Return(existingPayment, nil)

		mockChargeDao.On("Insert", mock.Anything).Return(&charge.Entity{}, nil)

		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, nil)
		mockOrderDao.On("Update", mock.Anything).Return(expectedOrder, nil)

		useCase := usecases.NewProcessPayment(mockPaymentDao, mockChargeDao, mockOrderDao)
		err := useCase.Execute(paymentID, processType, details)

		assert.NoError(t, err)
		mockPaymentDao.AssertExpectations(t)
		mockChargeDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})

	t.Run("should return error when payment not found", func(t *testing.T) {
		mockPaymentDao := new(testhelpers.MockPaymentDao)
		mockChargeDao := new(testhelpers.MockChargeDao)
		mockOrderDao := new(testhelpers.MockOrderDao)

		mockPaymentDao.On("FindById", paymentID).Return(nil, assert.AnError)

		useCase := usecases.NewProcessPayment(mockPaymentDao, mockChargeDao, mockOrderDao)
		err := useCase.Execute(paymentID, processType, details)

		assert.Error(t, err)
		mockPaymentDao.AssertExpectations(t)
		mockChargeDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})

	t.Run("should process payment successfully", func(t *testing.T) {
		mockPaymentDao := new(testhelpers.MockPaymentDao)
		mockChargeDao := new(testhelpers.MockChargeDao)
		mockOrderDao := new(testhelpers.MockOrderDao)
		var existingPayments []payment.Entity

		mockPaymentDao.On("FindById", paymentID).Return(existingPayment, nil)
		mockPaymentDao.On("FindByOrderId", mock.Anything, mock.Anything).Return(existingPayments, nil)
		mockPaymentDao.On("Update", mock.Anything).Return(existingPayment, nil)

		mockChargeDao.On("Insert", mock.Anything).Return(&charge.Entity{}, nil)

		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, nil)
		mockOrderDao.On("Update", mock.Anything).Return(expectedOrder, nil)

		useCase := usecases.NewProcessPayment(mockPaymentDao, mockChargeDao, mockOrderDao)
		err := useCase.Execute(paymentID, processType, details)

		assert.NoError(t, err)
		mockPaymentDao.AssertExpectations(t)
		mockChargeDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})

	t.Run("should throw error when payment update fails", func(t *testing.T) {
		mockPaymentDao := new(testhelpers.MockPaymentDao)
		mockChargeDao := new(testhelpers.MockChargeDao)
		mockOrderDao := new(testhelpers.MockOrderDao)
		var existingPayments []payment.Entity

		mockPaymentDao.On("FindById", paymentID).Return(existingPayment, nil)
		mockPaymentDao.On("FindByOrderId", mock.Anything, mock.Anything).Return(existingPayments, nil)
		mockPaymentDao.On("Update", mock.Anything).Return(existingPayment, assert.AnError)

		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, nil)

		useCase := usecases.NewProcessPayment(mockPaymentDao, mockChargeDao, mockOrderDao)
		err := useCase.Execute(paymentID, processType, details)

		assert.Error(t, err)
		mockPaymentDao.AssertExpectations(t)
		mockChargeDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})

	t.Run("should throw error when order update fails", func(t *testing.T) {
		mockPaymentDao := new(testhelpers.MockPaymentDao)
		mockChargeDao := new(testhelpers.MockChargeDao)
		mockOrderDao := new(testhelpers.MockOrderDao)
		var existingPayments []payment.Entity

		mockPaymentDao.On("FindById", paymentID).Return(existingPayment, nil)
		mockPaymentDao.On("FindByOrderId", mock.Anything, mock.Anything).Return(existingPayments, nil)
		mockPaymentDao.On("Update", mock.Anything).Return(existingPayment, nil)

		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, nil)
		mockOrderDao.On("Update", mock.Anything).Return(expectedOrder, assert.AnError)

		useCase := usecases.NewProcessPayment(mockPaymentDao, mockChargeDao, mockOrderDao)
		err := useCase.Execute(paymentID, processType, details)

		assert.Error(t, err)
		mockPaymentDao.AssertExpectations(t)
		mockChargeDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})

	t.Run("should throw error when charge insert fails", func(t *testing.T) {
		mockPaymentDao := new(testhelpers.MockPaymentDao)
		mockChargeDao := new(testhelpers.MockChargeDao)
		mockOrderDao := new(testhelpers.MockOrderDao)
		var existingPayments []payment.Entity

		mockPaymentDao.On("FindById", paymentID).Return(existingPayment, nil)
		mockPaymentDao.On("FindByOrderId", mock.Anything, mock.Anything).Return(existingPayments, nil)
		mockPaymentDao.On("Update", mock.Anything).Return(existingPayment, nil)

		mockChargeDao.On("Insert", mock.Anything).Return(&charge.Entity{}, assert.AnError)

		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, nil)
		mockOrderDao.On("Update", mock.Anything).Return(expectedOrder, nil)

		useCase := usecases.NewProcessPayment(mockPaymentDao, mockChargeDao, mockOrderDao)
		err := useCase.Execute(paymentID, processType, details)

		assert.Error(t, err)
		mockPaymentDao.AssertExpectations(t)
		mockChargeDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})

	t.Run("should return error when charge creation fails", func(t *testing.T) {
		mockPaymentDao := new(testhelpers.MockPaymentDao)
		mockChargeDao := new(testhelpers.MockChargeDao)
		mockOrderDao := new(testhelpers.MockOrderDao)
		var existingPayments []payment.Entity

		mockPaymentDao.On("FindById", paymentID).Return(existingPayment, nil)
		mockPaymentDao.On("FindByOrderId", mock.Anything, mock.Anything).Return(existingPayments, assert.AnError)

		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, nil)

		useCase := usecases.NewProcessPayment(mockPaymentDao, mockChargeDao, mockOrderDao)
		err := useCase.Execute(paymentID, processType, details)

		assert.Error(t, err)
		mockPaymentDao.AssertExpectations(t)
		mockChargeDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})

	t.Run("should return error when find order fails", func(t *testing.T) {
		mockPaymentDao := new(testhelpers.MockPaymentDao)
		mockChargeDao := new(testhelpers.MockChargeDao)
		mockOrderDao := new(testhelpers.MockOrderDao)

		mockPaymentDao.On("FindById", paymentID).Return(existingPayment, nil)

		mockOrderDao.On("FindById", mock.Anything).Return(expectedOrder, assert.AnError)

		useCase := usecases.NewProcessPayment(mockPaymentDao, mockChargeDao, mockOrderDao)
		err := useCase.Execute(paymentID, processType, details)

		assert.Error(t, err)
		mockPaymentDao.AssertExpectations(t)
		mockChargeDao.AssertExpectations(t)
		mockOrderDao.AssertExpectations(t)
	})
}
