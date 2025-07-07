package service

import (
	"database/sql"
	"fmt"

	postgres_db "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/internal/db"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/internal/parser"
)

type Service struct {
	Db postgres_db.DbService
}

type ServiceManager interface {
	ParserItem(link string) (string, string, float32, error)
	SelectItem(link string) (string, float32, error)
	UpdateItem(price float32, link string) error
	InsertItem(userId, link, name string, price float32) error
	SelectAllItems(userId string) (*sql.Rows, error)
}

func NewService(db postgres_db.DbService) *Service {
	return &Service{
		Db: db,
	}
}

func (s *Service) ParserItem(link string) (string, string, float32, error) {
	status := "Товар в наличии"
	name, price, err := parser.Parser(link)
	if err != nil {
		if err.Error() == "Товара нет в наличии" {
			status = err.Error()
		} else {
			return "", "", 0, fmt.Errorf("cannot parse this link: %v", err)
		}
	}

	return name, status, price, nil
}

func (s *Service) SelectItem(link string) (string, float32, error) {

	return s.Db.SelectItemFromDB(link)

}

func (s *Service) UpdateItem(price float32, link string) error {

	return s.Db.UpdateItemFromDB(price, link)

}

func (s *Service) InsertItem(userId, link, name string, price float32) error {

	return s.Db.InsertItemFromDB(userId, link, name, price)

}

func (s *Service) SelectAllItems(userId string) (*sql.Rows, error) {

	return s.Db.SelectAllItemsFromDB(userId)

}
