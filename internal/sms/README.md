# SMS Service

A robust and fault-tolerant SMS sending service built with Go, following
Domain-Driven Design (DDD) principles. This service provides reliable SMS
delivery through an external API, with features like exponential backoff
for retries, circuit breaker for failure handling, and asynchronous
processing.

## Features

- **Reliable SMS Delivery**: Sends SMS messages using an external API
  with configurable retry attempts to handle transient failures.
- **Circuit Breaker**: Uses the `gobreaker` library to prevent
  overwhelming the external API during sustained failures, ensuring
  system resilience.
- **Exponential Backoff**: Implements retry logic with the `backoff`
  library to handle temporary API issues.
- **Domain-Specific Error Handling**: Maps infrastructure errors to
  domain errors (e.g., `ErrServiceTimeout`, `ErrExternalServiceFailure`)
  for consistent error management.
- **Asynchronous Processing**: Supports non-blocking SMS sending with
  callback-based error handling via `AsyncSMSService`.
- **Context-Aware**: Respects client-provided context deadlines and
  enforces service-level timeouts.
- **Comprehensive Testing**: Includes unit tests covering success,
  failure, timeout, and circuit breaker scenarios.

## Architecture

The service follows **Domain-Driven Design (DDD)** principles, with a clear
separation of concerns:
- **Domain Layer**: Defines the `SMS` entity and domain-specific errors
  (`ErrServiceTimeout`, `ErrExternalServiceFailure`, `ErrCircuitOpen`,
  `ErrMaxRetries`).
- **Application Layer**: Implements the `SMSService` and `AsyncSMSService`
  for business logic, retry handling, and circuit breaker integration.
- **Infrastructure Layer**: Provides a mock implementation of the external
  SMS API (`MockExternalSMSClient`) for testing.

## Installation

1. Ensure Go 1.18 or later is installed.
2. Clone the repository:
   ```bash
   git clone https://github.com/alishojaeiir/mahdaad.git
   cd mahdaad
   ```
3. Install dependencies:
   ```bash
   go get github.com/cenkalti/backoff/v4@v4.3.0
   go get github.com/sony/gobreaker@v1.0.0
   ```

## Usage

### Running Tests
To run the unit tests:
```bash
go test -v ./internal/sms/infra
```

### Example
```go
package main

import (
	"context"
	"github.com/alishojaeiir/mahdaad/internal/sms/application"
	"github.com/alishojaeiir/mahdaad/internal/sms/domain"
	"github.com/alishojaeiir/mahdaad/internal/sms/infra"
	"time"
)

func main() {
	client := infra.NewExternalSMSClient(0.0)
	service := application.NewSMSService(client, 3, 2*time.Second)
	sms := domain.SMS{To: "+989123456789", Message: "Hello, World!"}

	err := service.SendSMS(context.Background(), sms)
	if err != nil {
		println("Error:", err.Error())
		return
	}
	println("SMS sent successfully")
}
```

## Configuration

- **MaxRetries**: Number of retry attempts for transient failures (default: 3).
- **Timeout**: Service-level timeout for each SMS request (default: 2 seconds).
- **Circuit Breaker**:
    - `MaxRequests`: 2 (in half-open state).
    - `Interval`: 60 seconds (reset failure counts).
    - `Timeout`: 30 seconds (time to transition to half-open state).
    - `ReadyToTrip`: Trips if 5+ requests have a failure ratio ≥ 60%.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue
for bugs, improvements, or new features.
