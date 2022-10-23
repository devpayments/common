package contracts

type Wallet interface {
	Charge()
	Fund()
	Hold()
}
