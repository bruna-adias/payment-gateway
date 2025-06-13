package payment

type Dao interface {
	FindById(id int64) (*Entity, error)
	FindByOrderId(id int64) ([]Entity, error)
	Insert(payment *Entity) (*Entity, error)
	Update(pay *Entity) (*Entity, error)
}
