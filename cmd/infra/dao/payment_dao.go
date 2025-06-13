package dao

import (
	"payment-gateway/cmd/domain/payment"
	"payment-gateway/cmd/infra/db"
	"time"
)

type PaymentModel struct {
	Id        int64
	OrderID   int64
	Amount    float64
	Status    string
	Type      string
	Details   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PaymentDao struct {
	db db.Client
}

func NewPaymentDao(db db.Client) *PaymentDao {
	return &PaymentDao{db: db}
}

func (p *PaymentDao) Insert(pay *payment.Entity) (*payment.Entity, error) {
	query := `INSERT INTO payments 
		(order_id, status, payment_type, amount, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	res, err := p.db.Exec(query,
		pay.OrderID(),
		pay.Status(),
		pay.Type(),
		pay.Amount(),
		pay.CreatedAt().Format("2006-01-02 15:04:05"),
		pay.UpdatedAt().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	pay.SetId(id)

	return pay, nil
}

func (p *PaymentDao) FindById(id int64) (*payment.Entity, error) {
	query := `SELECT id, order_id, status, payment_type, created_at, updated_at, IFNULL(details, '') as details, amount FROM payments WHERE id = ?`

	var pay PaymentModel

	row, err := p.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		err := row.Scan(&pay.Id, &pay.OrderID, &pay.Status, &pay.Type, &pay.CreatedAt, &pay.UpdatedAt, &pay.Details, &pay.Amount)
		if err != nil {
			return nil, err
		}
	}

	paymentEntity := payment.NewPaymentBuilder().WithId(pay.Id).
		WithOrderId(pay.OrderID).
		WithStatus(pay.Status).
		WithType(pay.Type).
		WithCreatedAt(pay.CreatedAt).
		WithDetails(pay.Details).
		WithAmount(pay.Amount).
		Build()

	return paymentEntity, nil
}

func (p *PaymentDao) FindByOrderId(id int64) ([]payment.Entity, error) {
	query := `SELECT id, order_id, status, payment_type, created_at, updated_at, IFNULL(details, '') as details, amount FROM payments WHERE order_id = ?`

	var payments []payment.Entity
	row, err := p.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var pay PaymentModel
		err := row.Scan(&pay.Id, &pay.OrderID, &pay.Status, &pay.Type, &pay.CreatedAt, &pay.UpdatedAt, &pay.Details, &pay.Amount)
		if err != nil {
			return nil, err
		}

		paymentEntity := payment.NewPaymentBuilder().WithId(pay.Id).
			WithOrderId(pay.OrderID).
			WithStatus(pay.Status).
			WithType(pay.Type).
			WithCreatedAt(pay.CreatedAt).
			WithDetails(pay.Details).
			WithAmount(pay.Amount).
			Build()

		payments = append(payments, *paymentEntity)
	}

	return payments, nil
}

func (p *PaymentDao) Update(pay *payment.Entity) (*payment.Entity, error) {
	query := `UPDATE payments 
		SET status = ?, updated_at = ?
		WHERE id = ?`

	_, err := p.db.Exec(query,
		pay.Status(),
		pay.UpdatedAt(),
		pay.Id(),
	)
	if err != nil {
		return nil, err
	}

	return pay, nil
}
