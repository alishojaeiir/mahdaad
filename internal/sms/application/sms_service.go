package application

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/alishojaeiir/mahdaad/internal/sms/domain"
	"github.com/cenkalti/backoff/v4"
	"github.com/sony/gobreaker"
)

type SMSService interface {
	SendSMS(ctx context.Context, sms domain.SMS) error
}

type ExternalSMSClient interface {
	Send(ctx context.Context, sms domain.SMS) error
}

type smsService struct {
	client         ExternalSMSClient
	circuitBreaker *gobreaker.CircuitBreaker
	maxRetries     int
	timeout        time.Duration
}

func NewSMSService(client ExternalSMSClient, maxRetries int, timeout time.Duration) SMSService {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "sms-service",
		MaxRequests: 2,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.6
		},
	})

	return &smsService{
		client:         client,
		circuitBreaker: cb,
		maxRetries:     maxRetries,
		timeout:        timeout,
	}
}

func (s *smsService) SendSMS(ctx context.Context, sms domain.SMS) error {
	sendCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.circuitBreaker.Execute(func() (interface{}, error) {
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = s.timeout

		err := backoff.Retry(func() error {
			if err := s.client.Send(sendCtx, sms); err != nil {
				if err.Error() == "external API failure" {
					return domain.ErrExternalServiceFailure
				}
				return err
			}
			return nil
		}, backoff.WithContext(b, sendCtx))

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				return nil, domain.ErrServiceTimeout
			}
			if b.GetElapsedTime() >= b.MaxElapsedTime {
				return nil, domain.ErrMaxRetries
			}
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return domain.ErrCircuitOpen
		}
		return err
	}

	return nil
}

type AsyncSMSService struct {
	smsService SMSService
	wg         sync.WaitGroup
}

func NewAsyncSMSService(smsService SMSService) *AsyncSMSService {
	return &AsyncSMSService{
		smsService: smsService,
	}
}

func (a *AsyncSMSService) SendSMSAsync(ctx context.Context, sms domain.SMS, callback func(error)) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		err := a.smsService.SendSMS(ctx, sms)
		if callback != nil {
			callback(err)
		}
	}()
}

func (a *AsyncSMSService) Wait() {
	a.wg.Wait()
}
