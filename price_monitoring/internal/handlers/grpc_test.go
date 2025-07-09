package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	proto "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/proto"
)

type ServiceManager interface {
	ParserItem(link string) (string, string, float32, error)
	SelectItem(link string) (string, float32, error)
	UpdateItem(price float32, link string) error
	InsertItem(userId, link, name string, price float32) error
	SelectAllItems(userId string) (*sql.Rows, error)
}

type MockService struct {
	ParserItemFunc     func(link string) (string, string, float32, error)
	SelectItemFunc     func(link string) (string, float32, error)
	UpdateItemFunc     func(price float32, link string) error
	InsertItemFunc     func(userId, link, name string, price float32) error
	SelectAllItemsFunc func(userId string) (*sql.Rows, error)
}

func (m MockService) SelectItem(link string) (string, float32, error) {

	return m.SelectItemFunc(link)

}

func (m MockService) UpdateItem(price float32, link string) error {

	return m.UpdateItemFunc(price, link)

}

func (m MockService) InsertItem(userId, link, name string, price float32) error {

	return m.InsertItemFunc(userId, link, name, price)

}

func (m MockService) SelectAllItems(userId string) (*sql.Rows, error) {

	return m.SelectAllItemsFunc(userId)

}

func (m MockService) ParserItem(link string) (string, string, float32, error) {

	return m.ParserItemFunc(link)
}

func TestGetItem(t *testing.T) {
	t.Run("func SelectItem return sql.ErrNoRow", func(t *testing.T) {
		mock := MockService{
			SelectItemFunc: func(link string) (string, float32, error) {
				return "", 0, sql.ErrNoRows
			},
			ParserItemFunc: func(link string) (string, string, float32, error) {
				return "TestItem", "TestStatus", 100.0, nil
			},
			InsertItemFunc: func(userId, link, name string, price float32) error {
				return nil
			},
		}

		req := &proto.GetItemRequest{
			UserId: "123",
			Link:   "TestLink.ru",
		}

		handler := NewHandler(mock)

		resp, err := handler.GetItem(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected err %v", err)
		}

		if resp.Item.StartPrice != 100.0 || resp.Item.Status != "TestStatus" || resp.Item.Name != "TestItem" {
			t.Errorf("expected StartPrice=%f Status=%s, Name=%s, got StartPrice=%f Status=%s, Name=%s",
				100.0, "TestStatus", "TestItem", resp.Item.StartPrice, resp.Item.Status, resp.Item.Name)
		}
	})

	t.Run("func SelectItem return item", func(t *testing.T) {
		mock := MockService{
			SelectItemFunc: func(link string) (string, float32, error) {
				return "TestItem", 100.0, nil
			},
			ParserItemFunc: func(link string) (string, string, float32, error) {
				return "TestItem", "TestStatus", 130.0, nil
			},
			UpdateItemFunc: func(price float32, link string) error {
				return nil
			},
		}

		req := &proto.GetItemRequest{
			UserId: "123",
			Link:   "TestLink.ru",
		}

		handler := NewHandler(mock)

		resp, err := handler.GetItem(context.Background(), req)
		if err != nil {
			t.Errorf("unexpected err %v", err)
		}

		if resp.Item.StartPrice != 100.0 || resp.Item.Status != "TestStatus" || resp.Item.Name != "TestItem" || resp.Item.DiffPrice != 30 {
			t.Errorf("expected StartPrice=%f Status=%s, Name=%s Different=%f, got StartPrice=%f Status=%s, Name=%s Different=%f",
				100.0, "TestStatus", "TestItem", 30.0, resp.Item.StartPrice, resp.Item.Status, resp.Item.Name, resp.Item.DiffPrice)
		}
	})

	t.Run("func InsertItem return error", func(t *testing.T) {
		mock := MockService{
			SelectItemFunc: func(link string) (string, float32, error) {
				return "", 0, sql.ErrNoRows
			},
			ParserItemFunc: func(link string) (string, string, float32, error) {
				return "TestItem", "TestStatus", 130.0, nil
			},
			InsertItemFunc: func(userId, link, name string, price float32) error {
				return fmt.Errorf("cannot insert item")
			},
		}

		req := &proto.GetItemRequest{
			UserId: "123",
			Link:   "TestLink.ru",
		}

		handler := NewHandler(mock)

		resp, err := handler.GetItem(context.Background(), req)
		if err.Error() != "cannot add new item Error: cannot insert item" {
			t.Errorf("unexpected err %v, expected %v", err, fmt.Errorf("cannot add new item Error: cannot insert item"))
		}

		if resp != nil {
			t.Errorf("resp should be nil")
		}
	})

	t.Run("func UpdateItem return error", func(t *testing.T) {
		mock := MockService{
			SelectItemFunc: func(link string) (string, float32, error) {
				return "TestItem", 100.0, nil
			},
			ParserItemFunc: func(link string) (string, string, float32, error) {
				return "TestItem", "TestStatus", 130.0, nil
			},
			UpdateItemFunc: func(price float32, link string) error {
				return fmt.Errorf("cannot update item")
			},
		}

		req := &proto.GetItemRequest{
			UserId: "123",
			Link:   "TestLink.ru",
		}

		handler := NewHandler(mock)

		resp, err := handler.GetItem(context.Background(), req)
		if err.Error() != "cannot update current_price Error: cannot update item" {
			t.Errorf("unexpected err %v, expected %v", err, fmt.Errorf("cannot update current_price Error: cannot update item"))
		}

		if resp != nil {
			t.Errorf("resp should be nil")
		}
	})

	t.Run("func ParserItem return error", func(t *testing.T) {
		mock := MockService{
			SelectItemFunc: func(link string) (string, float32, error) {
				return "", 0, sql.ErrNoRows
			},
			ParserItemFunc: func(link string) (string, string, float32, error) {
				return "", "", 0, fmt.Errorf("cannot parse item")
			},
		}

		req := &proto.GetItemRequest{
			UserId: "123",
			Link:   "TestLink.ru",
		}

		handler := NewHandler(mock)

		resp, err := handler.GetItem(context.Background(), req)
		if err.Error() != "cannot parse item" {
			t.Errorf("unexpected err %v, expected %v", err, fmt.Errorf("cannot parse item"))
		}

		if resp != nil {
			t.Errorf("resp should be nil")
		}
	})

}

func TestGetAllItems(t *testing.T) {

	t.Run("SelectAllItems return error", func(t *testing.T) {
		mock := MockService{
			SelectAllItemsFunc: func(userId string) (*sql.Rows, error) {
				return nil, fmt.Errorf("cannot select items")
			},
		}

		handler := NewHandler(mock)

		req := &proto.GetAllItemsRequest{
			UserId: "123",
		}

		resp, err := handler.GetAllItems(context.Background(), req)

		if err.Error() != "cannot get data from postgres Error: cannot select items" {
			t.Errorf("unexpected err:%v, expected %v", err, fmt.Errorf("cannot get data from postgres Error: cannot select items"))
		}

		if resp != nil {
			t.Errorf("resp should be nil")
		}

	})

	t.Run("SelectAllItems return sql.Rows", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"start_price", "link"}).
			AddRow(100.0, "http://example.com/item1").
			AddRow(200.0, "http://example.com/item2")

		mock.ExpectQuery("SELECT start_price, link FROM items WHERE user_id = ?").
			WithArgs("123").
			WillReturnRows(rows)

		mockS := MockService{
			SelectAllItemsFunc: func(userId string) (*sql.Rows, error) {
				return db.Query("SELECT start_price, link FROM items WHERE user_id = ?", userId)
			},
			ParserItemFunc: func(link string) (string, string, float32, error) {
				if link == "http://example.com/item1" {
					return "Item1", "Available", 110.0, nil
				}
				return "Item2", "Out of stock", 190.0, nil
			},
		}

		handler := NewHandler(mockS)

		req := &proto.GetAllItemsRequest{UserId: "123"}

		resp, err := handler.GetAllItems(context.Background(), req)
		assert.NoError(t, err)

		assert.Len(t, resp.Items, 2)
		assert.Equal(t, "Item1", resp.Items[0].Name)
		assert.Equal(t, float32(10.0), resp.Items[0].DiffPrice)
		assert.Equal(t, "Item2", resp.Items[1].Name)
		assert.Equal(t, float32(-10.0), resp.Items[1].DiffPrice)

		assert.NoError(t, mock.ExpectationsWereMet())

	})

}
