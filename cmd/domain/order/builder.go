package order

import "time"

type Builder struct {
	o *Entity
}

func NewOrderBuilder() *Builder {
	return &Builder{
		o: &Entity{
			createdAt: time.Now(),
		},
	}
}

func (b *Builder) WithId(id int64) *Builder {
	b.o.SetId(id)
	return b
}

func (b *Builder) WithStatus(status string) *Builder {
	b.o.SetStatus(status)
	return b
}

func (b *Builder) WithAmount(amount float64) *Builder {
	b.o.SetAmount(amount)
	return b
}

func (b *Builder) WithCreatedAt(createdAt time.Time) *Builder {
	b.o.SetCreatedAt(createdAt)
	return b
}

func (b *Builder) WithUpdatedAt(createdAt time.Time) *Builder {
	b.o.SetUpdateAt(createdAt)
	return b
}

func (b *Builder) Build() *Entity {
	return b.o
}
