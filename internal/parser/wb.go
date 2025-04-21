package parser

import (
	"fmt"

	"github.com/tebeka/selenium"
)

func WbParser(link string, driver selenium.WebDriver) (string, string, string, error) {

	err := driver.Get(link)
	if err != nil {
		return "", "", "", err
	}

	err = driver.Wait(func(wd selenium.WebDriver) (bool, error) {
		_, err := wd.FindElement(selenium.ByCSSSelector, "span.price-block__wallet-price.red-price")
		return err == nil, nil
	})
	if err != nil {
		return "", "", "", err
	}

	product_title, err := driver.FindElement(selenium.ByCSSSelector, "h1.product-page__title")
	if err != nil {
		return "", "", "", err
	}

	wallet_price, err := driver.FindElement(selenium.ByCSSSelector, "span.price-block__wallet-price.red-price")
	if err != nil {
		return "", "", "", err
	}

	price, err := driver.FindElement(selenium.ByCSSSelector, "ins.price-block__final-price.wallet")
	if err != nil {
		return "", "", "", err
	}
	title, _ := product_title.Text()
	priceWithCard, _ := wallet_price.Text()
	priceRegular, _ := price.Text()

	fmt.Println("!!!", priceRegular)

	return title, priceWithCard, priceRegular, nil
}
