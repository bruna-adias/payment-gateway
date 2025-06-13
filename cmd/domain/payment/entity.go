package payment

import (
	"time"
)

const (
	pendingStatus  = "pending"
	approvedStatus = "approved"
	reprovedStatus = "reproved"
)

type Entity struct {
	id     int64
	status string

	orderID     int64
	amount      float64
	paymentType string
	details     string

	createdAt time.Time
	updatedAt time.Time
}

func NewPayment(orderID int64, amount float64, paymentType string) *Entity {
	return &Entity{
		orderID:     orderID,
		amount:      amount,
		status:      pendingStatus,
		paymentType: paymentType,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

func (p *Entity) Process(processType string, details string) {
	p.details = details

	if processType == "Success" {
		p.approve()
	} else {
		p.reprove()
	}
}

func (p *Entity) approve() {
	p.status = approvedStatus
	p.updatedAt = time.Now()
}

func (p *Entity) reprove() {
	p.status = reprovedStatus
	p.updatedAt = time.Now()
}

func (p *Entity) Status() string {
	return p.status
}

func (p *Entity) OrderID() int64 {
	return p.orderID
}

func (p *Entity) Type() string {
	return p.paymentType
}

func (p *Entity) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Entity) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Entity) SetId(id int64) {
	p.id = id
}

func (p *Entity) Id() int64 {
	return p.id
}

func (p *Entity) Details() string {
	return p.details
}

func (p *Entity) Amount() float64 {
	return p.amount
}

func (p *Entity) SetOrderID(id int64) {
	p.orderID = id
	p.updatedAt = time.Now()
}

func (p *Entity) SetAmount(amount float64) {
	p.amount = amount
	p.updatedAt = time.Now()
}

func (p *Entity) SetStatus(status string) {
	p.status = status
	p.updatedAt = time.Now()
}

func (p *Entity) SetType(status string) {
	p.paymentType = status
	p.updatedAt = time.Now()
}

func (p *Entity) SetDetails(details string) {
	p.details = details
	p.updatedAt = time.Now()
}

func (p *Entity) SetCreatedAt(createdAt time.Time) {
	p.createdAt = createdAt
}

func (p *Entity) IsValid() bool {
	return p.status == approvedStatus
}
