package usecases

import (
	"payment-gateway/cmd/domain/charge"
	"payment-gateway/cmd/domain/order"
	"payment-gateway/cmd/domain/payment"
)

type GetCashout struct {
	orderDao   order.Dao
	paymentDao payment.Dao
	chargeDao  charge.Dao
}

type CashoutView struct {
	OrderId       int64   `json:"-"`
	CashedDebt    float64 `json:"cashed_debt"`
	RemainingDebt float64 `json:"remaining_debt"`
	Charges       float64 `json:"charges"`
	IsPaid        bool    `json:"is_paid"`
}

type Accountable interface {
	Amount() float64
}

func NewGetCashout(paymentDao payment.Dao, orderDao order.Dao, chargeDao charge.Dao) *GetCashout {
	return &GetCashout{
		paymentDao: paymentDao,
		orderDao:   orderDao,
		chargeDao:  chargeDao,
	}
}

func (c *GetCashout) Execute(orderId int64) (order.Entity, CashoutView, error) {
	or, err := c.orderDao.FindById(orderId)
	if err != nil {
		return order.Entity{}, CashoutView{}, err
	}

	paidAmount, err := GetPaidAmount(c.paymentDao, orderId)
	if err != nil {
		return order.Entity{}, CashoutView{}, err
	}

	charges, err := c.chargeDao.FindByOrderId(orderId)
	if err != nil {
		return order.Entity{}, CashoutView{}, err
	}

	var totalCharges float64
	for _, charge := range charges {
		totalCharges += charge.Amount()
	}

	return *or, CashoutView{
		OrderId:       orderId,
		CashedDebt:    paidAmount,
		RemainingDebt: or.Amount() - paidAmount,
		Charges:       totalCharges,
		IsPaid:        paidAmount >= or.Amount(),
	}, nil
}
