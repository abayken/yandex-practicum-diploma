package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

type DatabaseStorage struct {
	DB *pgx.Conn
}

func NewStorage(url string) *DatabaseStorage {
	conn, err := pgx.Connect(context.Background(), url)

	if err != nil {
		log.Fatal(err)
	}

	storage := &DatabaseStorage{DB: conn}
	storage.initTablesIfNeeded()

	return storage
}

func (storage *DatabaseStorage) initTablesIfNeeded() {
	_, err := storage.DB.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS USERS (ID SERIAL, LOGIN VARCHAR(100), PASSWORD VARCHAR(100));")
	_, err = storage.DB.Exec(context.Background(), "create table if not exists orders (id serial primary key, user_id integer, status varchar(100), number varchar(100) unique, added_at timestamp default now(), accrual integer);")

	if err != nil {
		log.Fatal(err)
	}
}
