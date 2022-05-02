package usecases

import (
	"strconv"

	"github.com/abayken/yandex-practicum-diploma/internal/custom_errors"
	"github.com/abayken/yandex-practicum-diploma/internal/helpers"
	"github.com/abayken/yandex-practicum-diploma/internal/repositories"
)

type WithdrawUseCase struct {
	OrdersRepo    repositories.OrdersRepository
	WithdrawsRepo repositories.WithdrawRepository
	UserUseCase   AuthUseCase
	Luhn          helpers.LuhnChecker
}

func (usecase WithdrawUseCase) Withdraw(userID int, orderNumber string, sum float32) error {
	balance, err := usecase.UserUseCase.GetBalance(userID)

	if err != nil {
		return err
	}

	if balance.Current < sum {
		return &custom_errors.InsufficientFundsError{}
	}

	order, err := strconv.Atoi(orderNumber)

	if err != nil {
		return err
	}

	if !usecase.Luhn.IsValid(order) {
		return &custom_errors.InvalidOrderNumber{}
	}

	err = usecase.WithdrawsRepo.Add(userID, sum, orderNumber)

	return err
}
