package dao

import (
	"payment-gateway/cmd/domain/charge"
	"payment-gateway/cmd/infra/db"
	"time"
)

type ChargeModel struct {
	Id        int64
	Amount    float64
	Category  string
	PaymentId int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChargeDao struct {
	db db.Client
}

func NewChargeDao(db db.Client) *ChargeDao {
	return &ChargeDao{db: db}
}

func (p *ChargeDao) Insert(c *charge.Entity) (*charge.Entity, error) {
	query := `INSERT INTO charges 
		( amount, category, payment_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`

	res, err := p.db.Exec(query,
		c.Amount(),
		c.Category(),
		c.PaymentId(),
		c.CreatedAt().Format("2006-01-02 15:04:05"),
		c.UpdatedAt().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	c.SetId(id)

	return c, nil
}

func (p *ChargeDao) FindById(id int64) (*charge.Entity, error) {
	query := `SELECT id, amount, category, payment_id, created_at, updated_at FROM charges WHERE id = ?`

	var model ChargeModel

	row, err := p.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		err := row.Scan(&model.Id, &model.Amount, &model.Category, &model.PaymentId, &model.CreatedAt, &model.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	chargeEntity := charge.NewChargeBuilder().
		WithId(model.Id).
		WithAmount(model.Amount).
		WithCategory(model.Category).
		WithPaymentId(model.PaymentId).
		WithCreatedAt(model.CreatedAt).
		WithUpdatedAt(model.UpdatedAt).
		Build()

	return chargeEntity, nil
}

func (p *ChargeDao) FindByOrderId(id int64) ([]charge.Entity, error) {
	query := `SELECT c.id, c.amount, c.category, c.payment_id, c.created_at, c.updated_at FROM charges c inner join payments p on c.payment_id = p.id where p.order_id = ?`

	var charges []charge.Entity
	row, err := p.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var model ChargeModel
		err := row.Scan(&model.Id, &model.Amount, &model.Category, &model.PaymentId, &model.CreatedAt, &model.UpdatedAt)
		if err != nil {
			return nil, err
		}

		chargeEntity := charge.NewChargeBuilder().
			WithId(model.Id).
			WithAmount(model.Amount).
			WithCategory(model.Category).
			WithPaymentId(model.PaymentId).
			WithCreatedAt(model.CreatedAt).
			WithUpdatedAt(model.UpdatedAt).
			Build()

		charges = append(charges, *chargeEntity)
	}

	return charges, nil
}
