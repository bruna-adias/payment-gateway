package testhelpers

import (
	"github.com/stretchr/testify/mock"
	"payment-gateway/cmd/domain/charge"
	"payment-gateway/cmd/domain/order"
	"payment-gateway/cmd/domain/payment"
)

type MockPaymentDao struct {
	mock.Mock
}

func (m *MockPaymentDao) FindById(id int64) (*payment.Entity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*payment.Entity), args.Error(1)
}

func (m *MockPaymentDao) FindByOrderId(id int64) ([]payment.Entity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]payment.Entity), args.Error(1)
}

func (m *MockPaymentDao) Insert(pay *payment.Entity) (*payment.Entity, error) {
	args := m.Called(pay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*payment.Entity), args.Error(1)
}

func (m *MockPaymentDao) Update(pay *payment.Entity) (*payment.Entity, error) {
	args := m.Called(pay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*payment.Entity), args.Error(1)
}

type MockOrderDao struct {
	mock.Mock
}

func (m *MockOrderDao) FindById(id int64) (*order.Entity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.Entity), args.Error(1)
}

func (m *MockOrderDao) Update(pay *order.Entity) (*order.Entity, error) {
	args := m.Called(pay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.Entity), args.Error(1)
}

type MockChargeDao struct {
	mock.Mock
}

func (m *MockChargeDao) Insert(pay *charge.Entity) (*charge.Entity, error) {
	args := m.Called(pay)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*charge.Entity), args.Error(1)
}

func (m *MockChargeDao) FindByOrderId(id int64) ([]charge.Entity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]charge.Entity), args.Error(1)
}
