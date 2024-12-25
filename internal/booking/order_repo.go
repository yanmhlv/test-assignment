package booking

import "sync"

var _ OrderRepository = (*InMemoryOrderRepository)(nil)

type OrderRepository interface {
	Create(order Order) error
}

type InMemoryOrderRepository struct {
	lock   sync.RWMutex
	orders []Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{}
}

func (r *InMemoryOrderRepository) Create(order Order) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.orders = append(r.orders, order)
	return nil
}
