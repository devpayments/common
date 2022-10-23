package contracts

import (
	"github.com/devpayments/common/model"
)

type PaymentDestination interface {
	Fund() (model.Transaction, error)
}
