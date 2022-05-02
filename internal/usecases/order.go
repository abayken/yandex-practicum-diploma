package usecases

import (
	"strconv"

	"github.com/abayken/yandex-practicum-diploma/internal/custom_errors"
	"github.com/abayken/yandex-practicum-diploma/internal/helpers"
	"github.com/abayken/yandex-practicum-diploma/internal/repositories"
	"github.com/jackc/pgx/v4"
)

type OrderUseCase struct {
	Repo repositories.OrdersRepository
	Luhn helpers.LuhnChecker
}

func (usecase OrderUseCase) Add(userID, orderNumber int) (bool, error) {
	if !usecase.Luhn.IsValid(orderNumber) {
		return false, &custom_errors.InvalidOrderNumber{}
	}

	order, err := usecase.Repo.GetOrder(userID, strconv.Itoa(orderNumber))

	if err == pgx.ErrNoRows {
		err := usecase.Repo.AddOrder(userID, strconv.Itoa(orderNumber), "NEW")

		if err != nil {
			return false, err
		} else {
			return true, err
		}
	}

	if err != nil {
		return false, err
	}

	return false, &custom_errors.OrderAlreadyAddedError{UserID: order.UserID}
}
