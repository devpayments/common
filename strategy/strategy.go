package strategy

import (
	"context"
	"github.com/devpayments/common/entity"
)

type Chargeable interface {
	Name() string
	InitiateCharge(ctx context.Context, source any, amount int64, currency entity.Currency) (entity.Transaction, error)
	GetChargeAuthorization(ctx context.Context, reference string) error
	AuthorizeCharge(ctx context.Context, reference string, authorizationData any) error
	CheckChargeStatus(ctx context.Context, reference string) error
	CompleteCharge(ctx context.Context, reference string) error
}

type Fundable interface {
	Name() string
	InitiateFunding(ctx context.Context, destination any, amount int64, currency entity.Currency) (entity.Transaction, error)
	CompleteFunding(ctx context.Context, reference string) error
}
