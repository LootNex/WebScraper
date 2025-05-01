package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/tebeka/selenium"
	proto "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price-monitoring/gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price-monitoring"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price-monitoring/internal/parser"
)

// "https://www.ozon.ru/product/zero-mileage-5w-30-maslo-motornoe-sinteticheskoe-1-l-1627408607/?at=Eqtk44V8ghrNGKRJTOPRG9LS0zoNA7UwDRn5mtNR6r8o"
// https://www.wildberries.ru/catalog/154859676/detail.aspx
type Service struct {
	proto.UnimplementedScraperServer
	Postg  *sql.DB
	Driver selenium.WebDriver
}

func (s *Service) GetItem(ctx context.Context, req *proto.GetItemRequest) (*proto.GetItemResponse, error) {

	var name string
	var start_price float32

	err := s.Postg.QueryRow("SELECT product_name, start_price FROM items WHERE user_id=$1", req.UserId).Scan(
		&name, &start_price)

	if err == sql.ErrNoRows {

		pars, err := parser.GetParser(req.Link, s.Driver)
		if err != nil {
			return nil, err
		}

		NewItem, err := pars.ParseLink()
		fmt.Println(NewItem)
		if err != nil {
			return nil, fmt.Errorf("cannot parse this link Error: %v", err)
		}

		id := uuid.New().String()
		price, _ := strconv.Atoi(NewItem.Price)
		_, err = s.Postg.Exec("INSERT INTO items VALUES ($1, $2, $3, $4, $5, $6, $7)",
			id, req.UserId, req.Link, NewItem.Name, float32(price), float32(price), time.Now())

		if err != nil {
			return nil, fmt.Errorf("cannot add new item Error: %v", err)
		}
		result := proto.ItemResponse{
			Name:         NewItem.Name,
			StartPrice:   float32(price),
			CurrentPrice: float32(price),
			DiffPrice:    0,
		}

		return &proto.GetItemResponse{
			Item: &result,
		}, nil

	} else if err != nil {
		return nil, err
	}

	pars, err := parser.GetParser(req.Link, s.Driver)
	if err != nil {
		return nil, err
	}

	price, err := pars.CheckPrice()
	if err != nil {
		return nil, err
	}

	_, err = s.Postg.Exec("UPDATE items SET current_price = $1 WHERE link = $2", price, req.Link)
	if err != nil {
		return nil, fmt.Errorf("cannot update current_price Error: %v", err)
	}

	result := proto.ItemResponse{
		Name:         name,
		StartPrice:   start_price,
		CurrentPrice: price,
		DiffPrice:    price - start_price,
	}

	return &proto.GetItemResponse{
		Item: &result,
	}, nil

}

func (s *Service) GetAllItems(ctx context.Context, req *proto.GetAllItemsRequest) (*proto.GetAllItemsResponse, error) {

	var items []*proto.ItemResponse

	rows, err := s.Postg.Query("SELECT start_price, link FROM items WHERE user_id = $1", req.UserId)
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

		pars, err := parser.GetParser(link, s.Driver)
		if err != nil {
			return nil, err
		}

		item, err := pars.ParseLink()
		if err != nil {
			return nil, fmt.Errorf("cannot parse this link Error: %v", err)
		}

		price, _ := strconv.Atoi(item.Price)
		items = append(items, &proto.ItemResponse{
			Name:         item.Name,
			StartPrice:   start_price,
			CurrentPrice: float32(price),
			DiffPrice:    float32(price) - start_price,
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
