package usecases

import (
	"strconv"

	"github.com/abayken/yandex-practicum-diploma/internal/errors"
	"github.com/abayken/yandex-practicum-diploma/internal/helpers"
	"github.com/abayken/yandex-practicum-diploma/internal/repositories"
	"github.com/jackc/pgx/v4"
)

type OrderUseCase struct {
	Repo repositories.OrdersRepository
	Luhn helpers.LuhnChecker
}

func (usecase OrderUseCase) Add(userID, orderNumber int) error {
	if !usecase.Luhn.IsValid(orderNumber) {
		return &errors.InvalidOrderNumber{}
	}

	order, err := usecase.Repo.GetOrder(userID, strconv.Itoa(orderNumber))

	if err == pgx.ErrNoRows {
		err := usecase.Repo.AddOrder(userID, strconv.Itoa(orderNumber), "NEW", 0)

		return err
	}

	if err != nil {
		return err
	}

	return &errors.OrderAlreadyAddedError{UserID: order.UserID}
}

func (usecase OrderUseCase) GetOrders(userID int) ([]repositories.Order, error) {
	return usecase.Repo.GetOrders(userID)
}
