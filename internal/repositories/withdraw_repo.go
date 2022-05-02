package repositories

import (
	"context"

	"github.com/abayken/yandex-practicum-diploma/internal/database"
)

type WithdrawRepository struct {
	Storage *database.DatabaseStorage
}

func (repo WithdrawRepository) GetTotalSumOfWithdrawn(userID int) (int, error) {
	db := repo.Storage.DB

	var sum int
	err := db.QueryRow(context.Background(), "SELECT COALESCE(SUM(SUM), 0) FROM TRANSACTIONS WHERE USER_ID = $1", userID).Scan(&sum)

	return sum, err
}

func (repo WithdrawRepository) Add(userID int, sum int, orderNumber string) error {
	_, err := repo.Storage.DB.Exec(context.Background(), "INSERT INTO TRANSACTIONS (USER_ID, ORDER_NUMBER, SUM) VALUES ($1, $2, $3)", userID, orderNumber, sum)

	return err
}
