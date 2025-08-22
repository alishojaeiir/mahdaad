package saga

import (
	"errors"
	"testing"
)

// Mock implementations for testing
type MockOrderService struct {
	CreateOrderFunc func() error
	CancelOrderFunc func() error
}

func (m *MockOrderService) CreateOrder() error {
	if m.CreateOrderFunc != nil {
		return m.CreateOrderFunc()
	}
	return nil
}

func (m *MockOrderService) CancelOrder() error {
	if m.CancelOrderFunc != nil {
		return m.CancelOrderFunc()
	}
	return nil
}

type MockInventoryService struct {
	DeductInventoryFunc  func() error
	AddBackInventoryFunc func() error
}

func (m *MockInventoryService) DeductInventory() error {
	if m.DeductInventoryFunc != nil {
		return m.DeductInventoryFunc()
	}
	return nil
}

func (m *MockInventoryService) AddBackInventory() error {
	if m.AddBackInventoryFunc != nil {
		return m.AddBackInventoryFunc()
	}
	return nil
}

type MockPaymentService struct {
	ProcessPaymentFunc func() error
	RefundPaymentFunc  func() error
}

func (m *MockPaymentService) ProcessPayment() error {
	if m.ProcessPaymentFunc != nil {
		return m.ProcessPaymentFunc()
	}
	return nil
}

func (m *MockPaymentService) RefundPayment() error {
	if m.RefundPaymentFunc != nil {
		return m.RefundPaymentFunc()
	}
	return nil
}

func TestSagaOrchestrator_Execute_Success(t *testing.T) {
	// Arrange
	orderSvc := &MockOrderService{}
	inventorySvc := &MockInventoryService{}
	paymentSvc := &MockPaymentService{}
	orchestrator := NewSagaOrchestrator(orderSvc, inventorySvc, paymentSvc)

	// Act
	err := orchestrator.Execute()

	// Assert
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestSagaOrchestrator_Execute_FailOnCreateOrder(t *testing.T) {
	// Arrange
	orderSvc := &MockOrderService{
		CreateOrderFunc: func() error { return ErrCreateOrderFailed },
	}
	inventorySvc := &MockInventoryService{}
	paymentSvc := &MockPaymentService{}
	orchestrator := NewSagaOrchestrator(orderSvc, inventorySvc, paymentSvc)

	// Act
	err := orchestrator.Execute()

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, ErrCreateOrderFailed) {
		t.Errorf("Expected 'order creation failed', got %v", err)
	}
}

func TestSagaOrchestrator_Execute_FailOnDeductInventory(t *testing.T) {
	// Arrange
	orderSvc := &MockOrderService{}
	inventorySvc := &MockInventoryService{
		DeductInventoryFunc: func() error { return ErrDeductInventoryFailed },
	}
	paymentSvc := &MockPaymentService{}
	orchestrator := NewSagaOrchestrator(orderSvc, inventorySvc, paymentSvc)

	// Act
	err := orchestrator.Execute()

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, ErrDeductInventoryFailed) {
		t.Errorf("Expected 'inventory deduction failed', got %v", err)
	}
}

func TestSagaOrchestrator_Execute_FailOnProcessPayment(t *testing.T) {
	// Arrange
	orderSvc := &MockOrderService{}
	inventorySvc := &MockInventoryService{}
	paymentSvc := &MockPaymentService{
		ProcessPaymentFunc: func() error { return ErrProcessPaymentFailed },
	}
	orchestrator := NewSagaOrchestrator(orderSvc, inventorySvc, paymentSvc)

	// Act
	err := orchestrator.Execute()

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, ErrProcessPaymentFailed) {
		t.Errorf("Expected 'payment processing failed', got %v", err)
	}
}
