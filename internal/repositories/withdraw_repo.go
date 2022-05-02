package repositories

import (
	"context"

	"github.com/abayken/yandex-practicum-diploma/internal/database"
)

type WithdrawRepository struct {
	Storage *database.DatabaseStorage
}

func (repo WithdrawRepository) GetTotalSumOfWithdrawn(userID int) (float32, error) {
	db := repo.Storage.DB

	var sum float32
	err := db.QueryRow(context.Background(), "SELECT COALESCE(SUM(SUM), 0) FROM TRANSACTIONS WHERE USER_ID = $1", userID).Scan(&sum)

	return sum, err
}

func (repo WithdrawRepository) Add(userID int, sum float32, orderNumber string) error {
	_, err := repo.Storage.DB.Exec(context.Background(), "INSERT INTO TRANSACTIONS (USER_ID, ORDER_NUMBER, SUM) VALUES ($1, $2, $3)", userID, orderNumber, sum)

	return err
}
