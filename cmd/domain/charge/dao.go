package charge

type Dao interface {
	Insert(charge *Entity) (*Entity, error)
	FindByOrderId(id int64) ([]Entity, error)
}
