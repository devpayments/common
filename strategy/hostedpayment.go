package strategy

type HostedPayment interface {
	GetAuthorizationURL()
	HandleWebhook()
	Confirm()
}
