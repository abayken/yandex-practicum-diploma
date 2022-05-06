package usecases

import (
	"github.com/abayken/yandex-practicum-diploma/internal/repositories"
)

type AccrualUseCase struct {
	OrdersRepository  repositories.OrdersRepository
	AccrualRepository repositories.AccrualRepository
}

func (usecase AccrualUseCase) ActualizeOrders(userID int) error {
	orders, err := usecase.OrdersRepository.GetNotFinishedOrders(userID)

	if err != nil {
		return err
	}

	updated := make(chan error)

	var updateError error

	for _, order := range orders {
		go usecase.update(updated, order.Number, userID)

		err := <-updated

		if err != nil {
			updateError = err
		}
	}

	return updateError
}

func (usecase AccrualUseCase) update(updated chan error, orderNumber string, userID int) {
	orderInfo, err := usecase.AccrualRepository.FetchOrderInfo(orderNumber)

	if err != nil {
		updated <- err

		return
	}

	err = usecase.OrdersRepository.Update(
		userID,
		orderInfo.Status,
		int(orderInfo.Accrual*100),
		orderInfo.Number,
	)

	updated <- err
}
