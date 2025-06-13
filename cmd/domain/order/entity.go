package order

import (
	"payment-gateway/cmd/domain/err"
	"payment-gateway/cmd/domain/payment"
	"time"
)

const (
	errPaymentExceedsDebt = "Payment exceeds debt"
	paidStatus            = "paid"
)

type Entity struct {
	id     int64
	status string
	amount float64

	createdAt time.Time
	updatedAt time.Time
}

func (o *Entity) Id() int64 {
	return o.id
}

func (o *Entity) Status() string {
	return o.status
}

func (o *Entity) Amount() float64 {
	return o.amount
}

func (o *Entity) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Entity) UpdatedAt() time.Time {
	return o.updatedAt
}

func (o *Entity) SetId(id int64) {
	o.id = id
}

func (o *Entity) SetStatus(status string) {
	o.status = status
	o.updatedAt = time.Now()
}

func (o *Entity) SetCreatedAt(at time.Time) {
	o.createdAt = at
}

func (o *Entity) SetUpdateAt(at time.Time) {
	o.updatedAt = at
}

func (o *Entity) SetAmount(amount float64) {
	o.amount = amount
	o.updatedAt = time.Now()
}

func (o *Entity) ProcessPayment(remainingDebt float64, pay payment.Entity) error {
	if pay.IsValid() {
		if remainingDebt < pay.Amount() {
			return exceptions.NewDomainError(errPaymentExceedsDebt)
		}

		if remainingDebt == pay.Amount() {
			o.paid()
		}
	}

	return nil
}

func (o *Entity) PreValidation(remainingDebt, amount float64) error {
	if remainingDebt < amount {
		return exceptions.NewDomainError(errPaymentExceedsDebt)
	}

	return nil
}

func (o *Entity) paid() {
	o.status = paidStatus
}
