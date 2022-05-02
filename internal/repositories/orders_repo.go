package repositories

import (
	"context"

	"github.com/abayken/yandex-practicum-diploma/internal/database"
)

type OrdersRepository struct {
	Storage *database.DatabaseStorage
}

type Order struct {
	UserID int
	Number string
}

func (repo OrdersRepository) GetOrder(userID int, orderNumber string) (Order, error) {
	db := repo.Storage.DB

	var order Order

	err := db.QueryRow(context.Background(), "SELECT USER_ID, NUMBER FROM ORDERS WHERE NUMBER = $1", orderNumber).Scan(&order.UserID, &order.Number)

	return order, err
}

func (repo OrdersRepository) AddOrder(userID int, orderNumber, status string) error {
	db := repo.Storage.DB

	_, err := db.Exec(context.Background(), "INSERT INTO ORDERS (USER_ID, NUMBER, STATUS) VALUES ($1, $2, $3)", userID, orderNumber, status)

	return err
}
