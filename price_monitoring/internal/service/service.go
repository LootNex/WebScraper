package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tebeka/selenium"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/internal/parser"
	proto "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/proto"
)

type Service struct {
	proto.UnimplementedScraperServer
	Postg  *sql.DB
	Driver selenium.WebDriver
}

func (s *Service) GetItem(ctx context.Context, req *proto.GetItemRequest) (*proto.GetItemResponse, error) {

	var name string
	var start_price float32

	err := s.Postg.QueryRow("SELECT product_name, start_price FROM auth.items WHERE link=$1", req.Link).Scan(
		&name, &start_price)

	if err == sql.ErrNoRows {

		status := "В наличие"
		name, price, err := parser.Parser(req.Link)
		if err != nil {
			if err.Error() == "Товара нет в наличии" {
				status = err.Error()
			} else {
				return nil, fmt.Errorf("cannot parse this link: %v", err)
			}
		}

		id := uuid.New().String()
		_, err = s.Postg.Exec("INSERT INTO auth.items VALUES ($1, $2, $3, $4, $5, $6, $7)",
			id, req.UserId, req.Link, name, price, price, time.Now())

		if err != nil {
			return nil, fmt.Errorf("cannot add new item Error: %v", err)
		}
		result := proto.ItemResponse{
			Name:         name,
			StartPrice:   price,
			CurrentPrice: price,
			DiffPrice:    0,
			Status:       status,
		}

		return &proto.GetItemResponse{
			Item: &result,
		}, nil

	} else if err != nil {
		return nil, err
	}

	status := "В наличие"
	name, price, err := parser.Parser(req.Link)
	if err != nil {
		if err.Error() == "Товара нет в наличии" {
			status = err.Error()
		} else {
			return nil, fmt.Errorf("cannot parse this link: %v", err)
		}
	}

	_, err = s.Postg.Exec("UPDATE auth.items SET current_price = $1 WHERE link = $2", price, req.Link)
	if err != nil {
		return nil, fmt.Errorf("cannot update current_price Error: %v", err)
	}

	result := proto.ItemResponse{
		Name:         name,
		StartPrice:   start_price,
		CurrentPrice: price,
		DiffPrice:    price - start_price,
		Status:       status,
	}

	return &proto.GetItemResponse{
		Item: &result,
	}, nil

}

func (s *Service) GetAllItems(ctx context.Context, req *proto.GetAllItemsRequest) (*proto.GetAllItemsResponse, error) {

	var items []*proto.ItemResponse

	rows, err := s.Postg.Query("SELECT start_price, link FROM auth.items WHERE user_id = $1", req.UserId)
	if err != nil {
		return nil, fmt.Errorf("cannot get data from postgres Error: %v", err)
	}
	var hasRows bool

	for rows.Next() {
		hasRows = true
		var link string
		var start_price float32

		if err := rows.Scan(&start_price, &link); err != nil {
			return nil, err
		}
		status := "В наличие"
		name, price, err := parser.Parser(link)
		if err != nil {
			if err.Error() == "Товара нет в наличии" {
				status = err.Error()
			} else {
				return nil, fmt.Errorf("cannot parse this link: %v", err)
			}
		}

		items = append(items, &proto.ItemResponse{
			Name:         name,
			StartPrice:   start_price,
			CurrentPrice: price,
			DiffPrice:    price - start_price,
			Status:       status,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot read strings Error: %v", err)
	}

	if !hasRows {
		return nil, fmt.Errorf("this user_id doesnot have links Error: %v", err)
	}

	return &proto.GetAllItemsResponse{Items: items}, nil

}
