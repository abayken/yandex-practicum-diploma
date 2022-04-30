package usecases

import (
	"github.com/abayken/yandex-practicum-diploma/internal/custom_errors"
	"github.com/abayken/yandex-practicum-diploma/internal/repositories"
	"github.com/brianvoe/sjwt"
	"github.com/jackc/pgx/v4"
)

type AuthUseCase struct {
	Repository repositories.AuthRepository
}

const jwtKey = "diploma"

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

		claims := sjwt.New()
		claims.Set("login", login)
		claims.Set("id", id)
		jwt := claims.Generate([]byte(jwtKey))

		return jwt, nil
	}
}

func (usecase AuthUseCase) Login(login, password string) (string, error) {
	user, err := usecase.Repository.Get(login, password)

	if err == nil {
		if user.Login == login && user.Password == password {
			claims := sjwt.New()
			claims.Set("login", user.Login)
			claims.Set("id", user.Id)
			jwt := claims.Generate([]byte(jwtKey))

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
