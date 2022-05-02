package usecases

import "github.com/abayken/yandex-practicum-diploma/internal/repositories"

type WithdrawUseCase struct {
	OrdersRepo    repositories.OrdersRepository
	WithdrawsRepo repositories.WithdrawRepository
}
