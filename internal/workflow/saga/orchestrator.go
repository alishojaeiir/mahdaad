package saga

import "errors"

var (
	ErrCreateOrderFailed     = errors.New("failed to create order")
	ErrDeductInventoryFailed = errors.New("failed to deduct inventory")
	ErrProcessPaymentFailed  = errors.New("failed to process payment")
)

// SagaOrchestrator coordinates the multi-step operation using the Saga pattern.
type SagaOrchestrator struct {
	orderService     OrderService
	inventoryService InventoryService
	paymentService   PaymentService
}

// NewSagaOrchestrator creates a new orchestrator with the given services.
func NewSagaOrchestrator(order OrderService, inventory InventoryService, payment PaymentService) *SagaOrchestrator {
	return &SagaOrchestrator{
		orderService:     order,
		inventoryService: inventory,
		paymentService:   payment,
	}
}

// Execute performs the multi-step flow with compensation only on failure.
// This is based on the Orchestration Saga pattern within Event-Driven Architecture (EDD),
// where a central coordinator manages the sequence and triggers compensating actions
// only if a step fails. Compensation is handled explicitly to avoid unnecessary execution.
func (s *SagaOrchestrator) Execute() error {
	var compensations []func() error // Stack of compensation actions

	if err := s.orderService.CreateOrder(); err != nil {
		return errors.Join(err, ErrCreateOrderFailed)
	}
	compensations = append(compensations, s.orderService.CancelOrder)

	if err := s.inventoryService.DeductInventory(); err != nil {
		s.executeCompensations(compensations)
		return errors.Join(err, ErrDeductInventoryFailed)
	}
	compensations = append(compensations, s.inventoryService.AddBackInventory)

	if err := s.paymentService.ProcessPayment(); err != nil {
		s.executeCompensations(compensations)
		return errors.Join(err, ErrProcessPaymentFailed)
	}

	return nil
}

// executeCompensations runs the compensation actions in reverse order.
func (s *SagaOrchestrator) executeCompensations(compensations []func() error) {
	for i := len(compensations) - 1; i >= 0; i-- {
		_ = compensations[i]() // Ignore compensation errors for simplicity
	}
}
