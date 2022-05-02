package usecases

import "github.com/abayken/yandex-practicum-diploma/internal/repositories"

type AccrualUseCase struct {
	OrdersRepository  repositories.OrdersRepository
	AccrualRepository repositories.AccrualRepository
}

func (usecase AccrualUseCase) ActualizeOrders(userID int) {
	orders, err := usecase.OrdersRepository.GetNotFinishedOrders(userID)

	if err != nil {
		return
	}

	for _, order := range orders {
		orderInfo, _ := usecase.AccrualRepository.FetchOrderInfo(order.Number)

		if orderInfo != nil {
			_ = usecase.OrdersRepository.Update(userID, orderInfo.Status, int(orderInfo.Accrual), orderInfo.Number)
		}
	}
}
