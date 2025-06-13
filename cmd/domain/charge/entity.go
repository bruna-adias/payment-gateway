package charge

import (
	"payment-gateway/cmd/domain/payment"
	"time"
)

var typeByCategory = map[string]string{
	"CreditCard": "financial_fee",
	"CashSlip":   "process_fee",
	"Cash":       "free",
}
var taxByType = map[string]float64{
	"financial_fee": 0.1,
	"process_fee":   0.2,
	"free":          0,
}

type Entity struct {
	amount    float64
	id        int64
	category  string
	paymentId int64

	createdAt time.Time
	updatedAt time.Time
}

func NewCharge(entity payment.Entity) (*Entity, bool) {
	if !entity.IsValid() {
		return nil, false
	}

	amount := getAmount(entity)
	category := getCategory(entity)
	paymentId := entity.Id()

	return &Entity{
		amount:    amount,
		category:  category,
		paymentId: paymentId,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, true
}

func getAmount(entity payment.Entity) float64 {
	category := getCategory(entity)

	return entity.Amount() * taxByType[category]
}

func getCategory(entity payment.Entity) string {
	category := entity.Type()
	t, ok := typeByCategory[category]
	if !ok {
		return "free"
	}

	return t
}

func (c *Entity) Amount() float64 {
	return c.amount
}

func (c *Entity) Id() int64 {
	return c.id
}

func (c *Entity) Category() string {
	return c.category
}

func (c *Entity) PaymentId() int64 {
	return c.paymentId
}

func (c *Entity) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Entity) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Entity) SetId(id int64) {
	c.id = id
}

func (c *Entity) SetCreatedAt(createdAt time.Time) {
	c.createdAt = createdAt
}

func (c *Entity) SetPaymentId(paymentId int64) {
	c.paymentId = paymentId
	c.updatedAt = time.Now()
}

func (c *Entity) SetCategory(category string) {
	c.category = category
	c.updatedAt = time.Now()
}

func (c *Entity) SetAmount(amount float64) {
	c.amount = amount
	c.updatedAt = time.Now()
}

func (c *Entity) SetUpdatedAt(at time.Time) {
	c.updatedAt = at
}
