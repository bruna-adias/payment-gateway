package usecases

import (
	"payment-gateway/cmd/domain/order"
	"payment-gateway/cmd/domain/payment"
)

type CreatePayment struct {
	paymentDao payment.Dao
	orderDao   order.Dao
}

func NewCreatePayment(paymentDao payment.Dao, orderDao order.Dao) *CreatePayment {
	return &CreatePayment{
		paymentDao: paymentDao,
		orderDao:   orderDao,
	}
}

func (c *CreatePayment) Execute(orderId int64, amount float64, status string) (*payment.Entity, error) {
	or, err := c.orderDao.FindById(orderId)
	if err != nil {
		return nil, err
	}
	pay := payment.NewPayment(orderId, amount, status)

	paidAmount, err := GetPaidAmount(c.paymentDao, or.Id())
	if err != nil {
		return nil, err
	}

	err = or.PreValidation(or.Amount()-paidAmount, pay.Amount())
	if err != nil {
		return nil, err
	}

	pay, err = c.paymentDao.Insert(pay)
	if err != nil {
		return nil, err
	}

	return pay, nil
}
