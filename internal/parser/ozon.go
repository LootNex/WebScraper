package parser

import (
	"fmt"
	"time"

	"github.com/tebeka/selenium"
)

func OzonParser(link string, driver selenium.WebDriver) (string, string, string, error) {

	err := driver.Get(link)
	if err != nil {
		return "", "", "", err
	}
	time.Sleep(6 * time.Second)
	text, err := driver.PageSource()
	fmt.Println(text, err)
	err = driver.Wait(func(wd selenium.WebDriver) (bool, error) {
		_, err := wd.FindElement(selenium.ByCSSSelector, "h1.mo5_28.tsHeadline550Medium")
		return err == nil, nil
	})
	if err != nil {
		return "", "", "", err
	}
	wallet_price, err := driver.FindElement(selenium.ByCSSSelector, "span.m4n_28.nm2_28")
	if err != nil {
		return "", "", "", err
	}

	price, err := driver.FindElement(selenium.ByCSSSelector, "span.mn9_28.m9n_28.om2_28")
	if err != nil {
		return "", "", "", err
	}

	product_title, err := driver.FindElement(selenium.ByCSSSelector, "h1.mo5_28.tsHeadline550Medium")
	if err != nil {
		return "", "", "", err
	}

	title, _ := product_title.Text()
	priceWithCard, _ := wallet_price.Text()
	priceRegular, _ := price.Text()

	return title, priceWithCard, priceRegular, nil

}
