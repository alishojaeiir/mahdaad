package infra

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/alishojaeiir/mahdaad/internal/sms/application"
	"github.com/alishojaeiir/mahdaad/internal/sms/domain"
)

func TestSMSService_SendSMS(t *testing.T) {
	tests := []struct {
		name        string
		failureRate float64
		expectedErr error
		ctx         context.Context
	}{
		{
			name:        "Successful SMS send",
			failureRate: 0.0,
			expectedErr: nil,
			ctx:         context.Background(),
		},
		{
			name:        "Failed SMS send with retries",
			failureRate: 1.0,
			expectedErr: domain.ErrExternalServiceFailure,
			ctx:         context.Background(),
		},
		{
			name:        "Timeout scenario",
			failureRate: 0.0,
			expectedErr: domain.ErrServiceTimeout,
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
				return ctx
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewExternalSMSClient(tt.failureRate)
			service := application.NewSMSService(client, 3, 2*time.Second)
			sms := domain.SMS{To: "+989123456789", Message: "Test message"}

			err := service.SendSMS(tt.ctx, sms)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestAsyncSMSService_SendSMSAsync(t *testing.T) {
	tests := []struct {
		name        string
		failureRate float64
		expectedErr error
		ctx         context.Context
	}{
		{
			name:        "Successful async SMS send",
			failureRate: 0.0,
			expectedErr: nil,
			ctx:         context.Background(),
		},
		{
			name:        "Failed async SMS send",
			failureRate: 1.0,
			expectedErr: domain.ErrExternalServiceFailure,
			ctx:         context.Background(),
		},
		{
			name:        "Async SMS send with context cancellation",
			failureRate: 0.0,
			expectedErr: domain.ErrServiceTimeout,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Immediately cancel the context
				return ctx
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewExternalSMSClient(tt.failureRate)
			service := application.NewSMSService(client, 3, 2*time.Second)
			asyncService := application.NewAsyncSMSService(service)

			var wg sync.WaitGroup
			var callbackErr error
			wg.Add(1)

			sms := domain.SMS{To: "+989123456789", Message: "Test message"}
			asyncService.SendSMSAsync(tt.ctx, sms, func(err error) {
				defer wg.Done()
				callbackErr = err
			})

			wg.Wait()
			if !errors.Is(callbackErr, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, callbackErr)
			}
		})
	}
}

func TestCircuitBreaker(t *testing.T) {
	tests := []struct {
		name        string
		failureRate float64
		expectedErr error
		actions     func(t *testing.T, service application.SMSService, sms domain.SMS, client application.ExternalSMSClient)
	}{
		{
			name:        "Circuit breaker opens after failures",
			failureRate: 1.0,
			expectedErr: domain.ErrCircuitOpen,
			actions: func(t *testing.T, service application.SMSService, sms domain.SMS, client application.ExternalSMSClient) {
				for i := 0; i < 5; i++ {
					err := service.SendSMS(context.Background(), sms)
					if err != nil && !errors.Is(err, domain.ErrExternalServiceFailure) && !errors.Is(err, domain.ErrCircuitOpen) {
						t.Errorf("Unexpected error during circuit breaker trigger: %v", err)
					}
					time.Sleep(100 * time.Millisecond)
				}
			},
		},
		{
			name:        "Circuit breaker recovers after timeout",
			failureRate: 1.0,
			expectedErr: nil,
			actions: func(t *testing.T, service application.SMSService, sms domain.SMS, client application.ExternalSMSClient) {
				for i := 0; i < 5; i++ {
					err := service.SendSMS(context.Background(), sms)
					if err != nil && !errors.Is(err, domain.ErrExternalServiceFailure) && !errors.Is(err, domain.ErrCircuitOpen) {
						t.Errorf("Unexpected error during circuit breaker trigger: %v", err)
					}
					time.Sleep(100 * time.Millisecond)
				}
				// Use type assertion to access mockExternalSMSClient and change FailureRate
				mockClient, ok := client.(*mockExternalSMSClient)
				if !ok {
					t.Fatal("Expected client to be *mockExternalSMSClient")
				}
				mockClient.failureRate = 0.0
				time.Sleep(31 * time.Second) // Wait for circuit breaker timeout (30s)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewExternalSMSClient(tt.failureRate)
			service := application.NewSMSService(client, 3, 2*time.Second)
			sms := domain.SMS{To: "+989123456789", Message: "Test message"}

			tt.actions(t, service, sms, client)
			err := service.SendSMS(context.Background(), sms)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
