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
