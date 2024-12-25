package booking

import "sync"

var _ OrderRepository = (*InMemoryOrderRepository)(nil)

type OrderRepository interface {
	Create(order Order) error
}

type InMemoryOrderRepository struct {
	lock   sync.Mutex
	orders []Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{}
}

func (r *InMemoryOrderRepository) Create(order Order) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// for events, use transactional outbox
	// in transaction, insert order and insert outbox event with payload

	r.orders = append(r.orders, order)
	return nil
}
