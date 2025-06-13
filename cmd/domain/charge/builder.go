package charge

import "time"

type Builder struct {
	pay *Entity
}

func NewChargeBuilder() *Builder {
	return &Builder{
		pay: &Entity{
			createdAt: time.Now(),
		},
	}
}

func (b *Builder) WithId(id int64) *Builder {
	b.pay.SetId(id)
	return b
}

func (b *Builder) WithAmount(amount float64) *Builder {
	b.pay.SetAmount(amount)
	return b
}

func (b *Builder) WithCategory(status string) *Builder {
	b.pay.SetCategory(status)
	return b
}

func (b *Builder) WithCreatedAt(createdAt time.Time) *Builder {
	b.pay.SetCreatedAt(createdAt)
	return b
}

func (b *Builder) WithUpdatedAt(createdAt time.Time) *Builder {
	b.pay.SetUpdatedAt(createdAt)
	return b
}

func (b *Builder) WithPaymentId(paymentId int64) *Builder {
	b.pay.SetPaymentId(paymentId)
	return b
}

func (b *Builder) Build() *Entity {
	return b.pay
}
