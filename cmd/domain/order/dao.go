package order

type Dao interface {
	FindById(id int64) (*Entity, error)
	Update(or *Entity) (*Entity, error)
}
