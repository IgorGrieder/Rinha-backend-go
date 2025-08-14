package ports

type PaymentService interface {
	ProcessPayment()
	GetAll()
}
