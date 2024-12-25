package booking

type Notifier interface {
	Notify(order Order) error
}

type NoopNotifier struct{}

func NewNoopNotifier() *NoopNotifier {
	return &NoopNotifier{}
}

func (n *NoopNotifier) Notify(order Order) error {
	return nil
}
