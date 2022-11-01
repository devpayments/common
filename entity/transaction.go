package entity

type Transaction struct {
	Type      string
	Status    string
	Amount    int64
	Currency  Currency
	Reference string
}
