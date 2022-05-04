package usecases

import (
	"strconv"

	"github.com/abayken/yandex-practicum-diploma/internal/custom_errors"
	"github.com/abayken/yandex-practicum-diploma/internal/helpers"
	"github.com/abayken/yandex-practicum-diploma/internal/repositories"
	"github.com/jackc/pgx/v4"
)

type WithdrawUseCase struct {
	OrdersRepo    repositories.OrdersRepository
	WithdrawsRepo repositories.WithdrawRepository
	UserUseCase   AuthUseCase
	Luhn          helpers.LuhnChecker
}

func (usecase WithdrawUseCase) Withdraw(userID int, orderNumber string, sum float32) error {
	orderInfo, err := usecase.OrdersRepo.GetOrder(userID, orderNumber)

	if err == pgx.ErrNoRows {
		order, err := strconv.Atoi(orderNumber)

		if err != nil {
			return err
		}

		if !usecase.Luhn.IsValid(order) {
			return &custom_errors.InvalidOrderNumber{}
		}

		err = usecase.OrdersRepo.AddOrder(userID, orderNumber, "NEW", 0)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		return &custom_errors.OrderAlreadyAddedError{UserID: orderInfo.UserID}
	}

	balance, err := usecase.UserUseCase.GetBalance(userID)

	if err != nil {
		return err
	}

	if balance.Current < sum {
		return &custom_errors.InsufficientFundsError{}
	}

	err = usecase.WithdrawsRepo.Add(userID, int(sum*100), orderNumber)

	return err
}

func (usecase WithdrawUseCase) Withdrawals(userID int) ([]repositories.Withdraw, error) {
	return usecase.WithdrawsRepo.GetWithdrawals(userID)
}
