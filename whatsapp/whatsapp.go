package whatsapp

type IClient interface {
	Connect(f func(s string)) (IClient, error)
	Send(phone string, message string) (err error)
	Closed()
}

func NewClient() IClient {
	return clientWhatsMeow()
}
