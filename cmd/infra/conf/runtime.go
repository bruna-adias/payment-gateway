package conf

import (
	"payment-gateway/cmd/infra"
	"payment-gateway/cmd/infra/dao"
	"payment-gateway/cmd/infra/db/mysql"
	"payment-gateway/cmd/infra/handler"
	"payment-gateway/cmd/usecases"
)

type Runtime struct {
	CreatePaymentHandler  handler.Handler
	ProcessPaymentHandler handler.Handler
	GetCashoutHandler     handler.Handler
}

func NewRuntime(configuration *infra.Configuration) *Runtime {
	// Create DB
	db, err := mysql.NewMySQLClient(configuration)
	if err != nil {
		panic(err)
	}

	// Create DAOs
	paymentDao := dao.NewPaymentDao(db)
	chargeDao := dao.NewChargeDao(db)
	orderDao := dao.NewOrderDao(db)

	// Create Use Cases
	createPayment := usecases.NewCreatePayment(paymentDao, orderDao)
	processPayment := usecases.NewProcessPayment(paymentDao, chargeDao, orderDao)
	getCashout := usecases.NewGetCashout(paymentDao, orderDao, chargeDao)

	// Create Handlers
	paymentHandler := handler.NewCreatePaymentHandler(createPayment)
	processPaymentHandler := handler.NewProcessPaymentHandler(processPayment)
	getCashoutHandler := handler.NewGetCashoutHandler(getCashout)

	return &Runtime{
		CreatePaymentHandler:  paymentHandler,
		ProcessPaymentHandler: processPaymentHandler,
		GetCashoutHandler:     getCashoutHandler,
	}
}
