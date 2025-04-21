package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tebeka/selenium"
)

type Item struct {
	Name  string
	Price string
}

type Parser interface {
	CheckPrice() (float32, error)
	ParseLink() (*Item, error)
}

type Ozon struct {
	link   string
	driver selenium.WebDriver
}
type Wb struct {
	link   string
	driver selenium.WebDriver
}

func GetParser(link string, driver selenium.WebDriver) (Parser, error) {

	switch {
	case strings.Contains(link, "ozon"):
		return &Ozon{
			link:   link,
			driver: driver,
		}, nil

	case strings.Contains(link, "wildberries"):
		return &Wb{
			link:   link,
			driver: driver,
		}, nil

	default:
		return nil, errors.New("unknown website")
	}

}

func (z *Ozon) CheckPrice() (float32, error) {

	item := &Item{}
	var err error

	_, _, item.Price, err = OzonParser(z.link, z.driver)
	if err != nil {
		return 0, fmt.Errorf("problems with parsing ozon Error: %v", err)
	}
	price, _ := strconv.Atoi(item.Price)
	return float32(price), nil
}

func (z *Ozon) ParseLink() (*Item, error) {

	item := &Item{}
	var err error

	item.Name, _, item.Price, err = OzonParser(z.link, z.driver)
	if err != nil {
		return nil, fmt.Errorf("problems with parsing ozon Error: %v", err)
	}
	return item, nil

}

func (w *Wb) CheckPrice() (float32, error) {

	item := &Item{}
	var err error

	_, _, item.Price, err = WbParser(w.link, w.driver)
	if err != nil {
		return 0, fmt.Errorf("problems with parsing wb Error: %v", err)
	}
	price, _ := strconv.Atoi(item.Price)
	return float32(price), nil

}

func (w *Wb) ParseLink() (*Item, error) {

	item := &Item{}
	var err error

	item.Name, _, item.Price, err = WbParser(w.link, w.driver)
	fmt.Println(&item.Price)
	if err != nil {
		return nil, fmt.Errorf("problems with parsing wb Error: %v", err)
	}
	return item, nil

}
