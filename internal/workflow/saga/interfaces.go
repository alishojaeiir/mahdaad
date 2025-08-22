package saga

// OrderService defines the interface for order operations.
type OrderService interface {
	CreateOrder() error
	CancelOrder() error // Compensation action
}

// InventoryService defines the interface for inventory operations.
type InventoryService interface {
	DeductInventory() error
	AddBackInventory() error // Compensation action
}

// PaymentService defines the interface for payment operations.
type PaymentService interface {
	ProcessPayment() error
	RefundPayment() error // Compensation action
}
