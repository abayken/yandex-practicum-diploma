package repositories

import (
	"context"

	"github.com/abayken/yandex-practicum-diploma/internal/database"
)

type AuthRepository struct {
	Storage *database.DatabaseStorage
}

type User struct {
	ID       int
	Login    string
	Password string
}

func (repo *AuthRepository) Exists(login string) (bool, error) {
	db := repo.Storage.DB

	var exists bool
	err := db.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM USERS WHERE LOGIN = $1);", login).Scan(&exists)

	if err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

func (repo *AuthRepository) Create(login, password string) (int, error) {
	db := repo.Storage.DB

	var id int
	err := db.QueryRow(context.Background(), "INSERT INTO USERS (LOGIN, PASSWORD) VALUES ($1, $2) RETURNING ID;", login, password).Scan(&id)

	return id, err
}

func (repo *AuthRepository) Get(login, password string) (User, error) {
	db := repo.Storage.DB

	var user User

	err := db.QueryRow(context.Background(), "SELECT ID, LOGIN, PASSWORD FROM USERS WHERE LOGIN = $1", login).Scan(&user.ID, &user.Login, &user.Password)

	return user, err
}
