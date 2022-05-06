package repositories

import (
	"context"

	"github.com/abayken/yandex-practicum-diploma/internal/database"
	"github.com/jackc/pgtype"
)

type WithdrawRepository struct {
	Storage *database.DatabaseStorage
}

type Withdraw struct {
	OrderNumber string
	Sum         int
	AddedAt     pgtype.Timestamp
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

func (repo WithdrawRepository) GetWithdrawals(userID int) ([]Withdraw, error) {
	rows, err := repo.Storage.DB.Query(context.Background(), "SELECT ORDER_NUMBER, SUM, ADDED_AT FROM TRANSACTIONS WHERE USER_ID = $1 ORDER BY ADDED_AT ASC;", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var withdrawals []Withdraw

	for rows.Next() {
		var withdraw Withdraw

		err := rows.Scan(&withdraw.OrderNumber, &withdraw.Sum, &withdraw.AddedAt)

		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdraw)
	}

	return withdrawals, nil
}
