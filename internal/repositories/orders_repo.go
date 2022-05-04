package repositories

import (
	"context"

	"github.com/abayken/yandex-practicum-diploma/internal/database"
	"github.com/jackc/pgtype"
)

type OrdersRepository struct {
	Storage *database.DatabaseStorage
}

type Order struct {
	UserID  int
	Number  string
	Status  string
	AddedAt pgtype.Timestamp
	Accrual int
}

func (repo OrdersRepository) GetOrder(userID int, orderNumber string) (Order, error) {
	db := repo.Storage.DB

	var order Order

	err := db.QueryRow(context.Background(), "SELECT USER_ID, NUMBER FROM ORDERS WHERE NUMBER = $1", orderNumber).Scan(&order.UserID, &order.Number)

	return order, err
}

func (repo OrdersRepository) AddOrder(userID int, orderNumber, status string, accrual int) error {
	db := repo.Storage.DB

	_, err := db.Exec(context.Background(), "INSERT INTO ORDERS (USER_ID, NUMBER, STATUS, ACCRUAL) VALUES ($1, $2, $3, $4)", userID, orderNumber, status, accrual)

	return err
}

func (repo OrdersRepository) GetOrders(userID int) ([]Order, error) {
	db := repo.Storage.DB

	rows, err := db.Query(context.Background(), "SELECT NUMBER, STATUS, ADDED_AT, ACCRUAL FROM ORDERS WHERE USER_ID = $1 ORDER BY ADDED_AT ASC;", userID)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	orders := []Order{}

	for rows.Next() {
		var order Order

		err := rows.Scan(&order.Number, &order.Status, &order.AddedAt, &order.Accrual)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (repo OrdersRepository) Update(status string, accrual int, number string) error {
	db := repo.Storage.DB

	_, err := db.Exec(context.Background(), "UPDATE ORDERS SET STATUS = $1, ACCRUAL = $2 WHERE NUMBER = $3", status, accrual, number)

	return err
}

func (repo OrdersRepository) GetNotFinishedOrders(userID int) ([]Order, error) {
	db := repo.Storage.DB

	rows, err := db.Query(context.Background(), "SELECT NUMBER FROM ORDERS WHERE USER_ID = $1 AND STATUS NOT IN ('INVALID', 'PROCESSED')", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order

		err := rows.Scan(&order.Number)

		if err == nil {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (repo OrdersRepository) GetAccrualSum(userID int) (int, error) {
	db := repo.Storage.DB

	var sum int
	err := db.QueryRow(context.Background(), "SELECT COALESCE(SUM(ACCRUAL), 0) FROM ORDERS WHERE USER_ID = $1 AND STATUS = 'PROCESSED';", userID).Scan(&sum)

	return sum, err
}
