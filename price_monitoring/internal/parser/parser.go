package parser

import (
	"os/exec"
	"strconv"
	"strings"
)

func Parser(link string) (string, float32, error) {
	cmd := exec.Command("/usr/bin/python3", "./scraper.py", link)

	output, err := cmd.Output()
	if err != nil {
		return "", 0, err
	}
	product_info := strings.Split(string(output[2:len(output)-2]), "'")
	name := product_info[0]
	product_price := product_info[1][2:]
	price, err := strconv.Atoi(product_price)
	if err != nil {
		return "", 0, err
	}

	return name, float32(price), nil
}
