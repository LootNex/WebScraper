package handlers

import (
	"context"
	"database/sql"
	"fmt"

	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/internal/service"
	proto "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/proto"
)

type Handler struct {
	proto.UnimplementedScraperServer
	Serv service.ServiceManager
}

func NewHandler(service service.ServiceManager) *Handler {
	return &Handler{
		Serv: service,
	}
}

func (s *Handler) GetItem(ctx context.Context, req *proto.GetItemRequest) (*proto.GetItemResponse, error) {

	_, start_price, err := s.Serv.SelectItem(req.Link)

	if err == sql.ErrNoRows {

		name, status, price, err := s.Serv.ParserItem(req.Link)
		if err != nil {
			return nil, err
		}

		err = s.Serv.InsertItem(req.UserId, req.Link, name, price)

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

	name, status, price, err := s.Serv.ParserItem(req.Link)
	if err != nil {
		return nil, err
	}

	err = s.Serv.UpdateItem(price, req.Link)
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

func (s *Handler) GetAllItems(ctx context.Context, req *proto.GetAllItemsRequest) (*proto.GetAllItemsResponse, error) {

	var items []*proto.ItemResponse

	rows, err := s.Serv.SelectAllItems(req.UserId)
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
		name, status, price, err := s.Serv.ParserItem(link)
		if err != nil {
			return nil, err
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
