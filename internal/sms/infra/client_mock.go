package infra

import (
	"context"
	"errors"
	"github.com/alishojaeiir/mahdaad/internal/sms/application"
	"github.com/alishojaeiir/mahdaad/internal/sms/domain"
	"math/rand"
)

type mockExternalSMSClient struct {
	failureRate float64
}

func NewExternalSMSClient(failureRate float64) application.ExternalSMSClient {
	return &mockExternalSMSClient{
		failureRate: failureRate,
	}
}

func (c *mockExternalSMSClient) Send(ctx context.Context, sms domain.SMS) error {
	select {
	case <-ctx.Done():
		return domain.ErrServiceTimeout
	default:
		if rand.Float64() < c.failureRate {
			return errors.New("external API failure")
		}
		return nil
	}
}
