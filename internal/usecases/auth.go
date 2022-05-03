package usecases

import (
	"github.com/abayken/yandex-practicum-diploma/internal/creds"
	"github.com/abayken/yandex-practicum-diploma/internal/custom_errors"
	"github.com/abayken/yandex-practicum-diploma/internal/repositories"
	"github.com/jackc/pgx/v4"
)

type AuthUseCase struct {
	Repository    repositories.AuthRepository
	Creds         creds.Creds
	OrdersRepo    repositories.OrdersRepository
	WithdrawsRepo repositories.WithdrawRepository
}

func (usecase AuthUseCase) Register(login, password string) (string, error) {
	exists, err := usecase.Repository.Exists(login)

	if err != nil {
		return "", err
	}

	if exists {
		return "", &custom_errors.AlreadyExistsUserError{}
	} else {
		id, err := usecase.Repository.Create(login, password)

		if err != nil {
			return "", err
		}

		jwt := usecase.Creds.BuildJWT(creds.UserModel{Login: login, Id: id})

		return jwt, nil
	}
}

func (usecase AuthUseCase) Login(login, password string) (string, error) {
	user, err := usecase.Repository.Get(login, password)

	if err == nil {
		if user.Login == login && user.Password == password {
			jwt := usecase.Creds.BuildJWT(creds.UserModel{Login: user.Login, Id: user.Id})

			return jwt, nil
		} else {
			return "", &custom_errors.InvalidCredentialsError{}
		}
	} else {
		if err == pgx.ErrNoRows {
			return "", &custom_errors.InvalidCredentialsError{}
		} else {
			return "", err
		}
	}
}

type Balance struct {
	Current        float32
	TotalWithdrawn float32
}

func (usecase AuthUseCase) GetBalance(userID int) (*Balance, error) {
	accrualSum, err := usecase.OrdersRepo.GetAccrualSum(userID)

	if err != nil {
		return nil, err
	}

	withdrawnSum, err := usecase.WithdrawsRepo.GetTotalSumOfWithdrawn(userID)

	if err != nil {
		return nil, err
	}

	return &Balance{
			Current:        float32(accrualSum-withdrawnSum) / 100,
			TotalWithdrawn: float32(withdrawnSum) / 100,
		},
		nil
}
