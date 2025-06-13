package usecases

import "payment-gateway/cmd/domain/payment"

func GetPaidAmount(dao payment.Dao, orderId int64) (float64, error) {
	payments, err := dao.FindByOrderId(orderId)
	if err != nil {
		return 0, err
	}

	var paidAmount float64
	for _, payment := range payments {
		if !payment.IsValid() {
			continue
		}
		paidAmount += payment.Amount()
	}

	return paidAmount, nil
}
