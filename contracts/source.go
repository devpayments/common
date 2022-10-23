package contracts

import (
	"github.com/devpayments/common/model"
)

type PaymentSource interface {
	Charge() (model.Transaction, error)
}
