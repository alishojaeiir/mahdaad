# Saga Orchestrator

This package implements a `SagaOrchestrator` using the **Orchestration Saga pattern** within an **Event-Driven Architecture (EDD)**. It coordinates multi-step operations (e.g., creating an order, deducting inventory, and processing payment) across independent services, with built-in compensation logic to handle failures.

## Overview
The `SagaOrchestrator` manages a sequence of steps, ensuring that if any step fails, previous actions are rolled back using compensating transactions. This approach is suitable for decoupled microservices or distributed systems where traditional ACID transactions are not feasible.

## Features
- Executes a multi-step flow (create order → deduct inventory → process payment).
- Provides compensation actions to undo previous steps on failure.
- Uses a central coordinator to manage the Saga pattern.
- Includes unit tests for success and failure scenarios.

## Error Handling
The package defines the following error constants:

`ErrCreateOrderFailed`: Indicates a failure in creating an order.
`ErrDeductInventoryFailed`: Indicates a failure in deducting inventory.
`ErrProcessPaymentFailed`: Indicates a failure in processing payment.

Errors are combined with these constants using errors.Join for detailed reporting.

## Testing
Run the tests using:
```bash
go test ./internal/workflow/saga
```
For verbose output:
```bash
go test -v ./internal/workflow/saga
```
