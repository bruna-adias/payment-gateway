package usecases

import (
	"payment-gateway/cmd/domain/charge"
	"payment-gateway/cmd/domain/order"
	"payment-gateway/cmd/domain/payment"
)

type ProcessPayment struct {
	paymentDao payment.Dao
	chargeDao  charge.Dao
	orderDao   order.Dao
}

func NewProcessPayment(paymentDao payment.Dao, chargeDao charge.Dao, orderDao order.Dao) *ProcessPayment {
	return &ProcessPayment{
		paymentDao: paymentDao,
		chargeDao:  chargeDao,
		orderDao:   orderDao,
	}
}

func (p *ProcessPayment) Execute(paymentID int64, processType string, details string) error {
	pay, err := p.paymentDao.FindById(paymentID)
	if err != nil {
		return err
	}

	or, err := p.orderDao.FindById(pay.OrderID())
	if err != nil {
		return err
	}

	paidAmount, err := GetPaidAmount(p.paymentDao, or.Id())
	if err != nil {
		return err
	}

	pay.Process(processType, details)
	err = or.ProcessPayment(or.Amount()-paidAmount, *pay)
	if err != nil {
		return err
	}

	_, err = p.paymentDao.Update(pay)
	if err != nil {
		return err
	}
	_, err = p.orderDao.Update(or)
	if err != nil {
		return err
	}

	newCharge, ok := charge.NewCharge(*pay)
	if ok {
		newCharge, err = p.chargeDao.Insert(newCharge)
		if err != nil {
			return err
		}
	}

	return nil
}
