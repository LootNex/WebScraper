package postgres_db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type DBConn struct {
	Conn *sql.DB
}

type DbService interface {
	SelectItemFromDB(link string) (string, float32, error)
	InsertItemFromDB(userId, link, name string, price float32) error
	UpdateItemFromDB(price float32, link string) error
	SelectAllItemsFromDB(userId string) (*sql.Rows, error)
}

func InitDB(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME string) (*DBConn, error) {

	str_conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	postg_conn, err := sql.Open("postgres", str_conn)

	if err != nil {
		return nil, err
	}

	err = postg_conn.Ping()
	if err != nil {
		return nil, err
	}

	return &DBConn{Conn: postg_conn}, nil

}

func (db *DBConn) SelectItemFromDB(link string) (string, float32, error) {

	var name string
	var start_price float32

	err := db.Conn.QueryRow("SELECT product_name, start_price FROM auth.items WHERE link=$1", link).Scan(
		&name, &start_price)

	return name, start_price, err

}

func (db *DBConn) InsertItemFromDB(userId, link, name string, price float32) error {

	id := uuid.New().String()

	_, err := db.Conn.Exec("INSERT INTO auth.items VALUES ($1, $2, $3, $4, $5, $6, $7)",
		id, userId, link, name, price, price, time.Now())

	return err
}

func (db *DBConn) UpdateItemFromDB(price float32, link string) error {

	_, err := db.Conn.Exec("UPDATE auth.items SET current_price = $1 WHERE link = $2", price, link)

	return err
}

func (db *DBConn) SelectAllItemsFromDB(userId string) (*sql.Rows, error) {

	rows, err := db.Conn.Query("SELECT start_price, link FROM auth.items WHERE user_id = $1", userId)

	return rows, err
}
