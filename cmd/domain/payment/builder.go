package payment

import "time"

type Builder struct {
	pay *Entity
}

func NewPaymentBuilder() *Builder {
	return &Builder{
		pay: &Entity{
			createdAt: time.Now(),
			status:    pendingStatus,
		},
	}
}

func (b *Builder) WithId(id int64) *Builder {
	b.pay.SetId(id)
	return b
}

func (b *Builder) WithOrderId(id int64) *Builder {
	b.pay.SetOrderID(id)
	return b
}

func (b *Builder) WithAmount(amount float64) *Builder {
	b.pay.SetAmount(amount)
	return b
}

func (b *Builder) WithStatus(status string) *Builder {
	b.pay.SetStatus(status)
	return b
}

func (b *Builder) WithType(status string) *Builder {
	b.pay.SetType(status)
	return b
}

func (b *Builder) WithDetails(details string) *Builder {
	b.pay.SetDetails(details)
	return b
}

func (b *Builder) WithCreatedAt(createdAt time.Time) *Builder {
	b.pay.SetCreatedAt(createdAt)
	return b
}

func (b *Builder) Build() *Entity {
	return b.pay
}
