package postgres_db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"os"
)

func InitDB() (*sql.DB, error) {

	str_conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
		
	postg_conn, err := sql.Open("postgres", str_conn)

	if err != nil {
		return nil, err
	}

	err = postg_conn.Ping()
	if err != nil {
		return nil, err
	}

	return postg_conn, nil

}
