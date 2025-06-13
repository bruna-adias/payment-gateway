package dao

import (
	"payment-gateway/cmd/domain/order"
	"payment-gateway/cmd/infra/db"
	"time"
)

type OrderModel struct {
	Id        int64
	Amount    float64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderDao struct {
	db db.Client
}

func NewOrderDao(db db.Client) *OrderDao {
	return &OrderDao{db: db}
}

func (p *OrderDao) FindById(id int64) (*order.Entity, error) {
	query := `SELECT id, status, amount, created_at, updated_at FROM orders WHERE id = ?`

	var pay OrderModel

	row, err := p.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		err := row.Scan(&pay.Id, &pay.Status, &pay.Amount, &pay.CreatedAt, &pay.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	orderEntity := order.NewOrderBuilder().WithId(pay.Id).
		WithStatus(pay.Status).
		WithAmount(pay.Amount).
		WithCreatedAt(pay.CreatedAt).
		WithUpdatedAt(pay.UpdatedAt).
		Build()

	return orderEntity, nil
}

func (p *OrderDao) Update(or *order.Entity) (*order.Entity, error) {
	query := `UPDATE orders 
		SET status = ?, updated_at = ?
		WHERE id = ?`

	_, err := p.db.Exec(query,
		or.Status(),
		or.UpdatedAt(),
		or.Id(),
	)
	if err != nil {
		return nil, err
	}

	return or, nil
}
