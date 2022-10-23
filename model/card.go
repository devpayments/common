package model

import "context"

type Card struct {
	FirstName        string
	LastName         string
	Number           string
	ExpiryMonth      int
	ExpiryYear       int
	StartMonth       int
	StartYear        int
	CVV              string
	IssueNumber      string
	Brand            string
	BillingAddress1  string
	BillingAddress2  string
	BillingCity      string
	BillingPostcode  string
	BillingState     string
	BillingCountry   string
	BillingPhone     string
	ShippingAddress1 string
	ShippingAddress2 string
	ShippingCity     string
	ShippingPostcode string
	ShippingState    string
	ShippingCountry  string
	ShippingPhone    string
	Company          string
	Email            string
}

func NewCard() Card {
	return Card{}
}

func (c *Card) ValidateCard(ctx context.Context) (error, bool) {
	return nil, true
}

func (c *Card) GetMaskedPan(ctx context.Context) string {
	return c.Number
}
