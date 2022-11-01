package entity

type Currency string
type PaymentType string

const (
	NigeriaNaira     = Currency("NGN")
	UKPounds         = Currency("GBP")
	USDollar         = Currency("USD")
	CanadianDollar   = Currency("CAD")
	KenyaShilling    = Currency("KES")
	GhanaCedis       = Currency("GHS")
	Euro             = Currency("EUR")
	USDT             = Currency("USDT")
	TanzaniaShilling = Currency("TZS")
	CentralFranc     = Currency("XAF")
	UgandaShilling   = Currency("UGX")
	WesternFranc     = Currency("XOF")
	RwandaFranc      = Currency("RWF")
)

const (
	CustomerPayment = PaymentType("customer_payment")
)
