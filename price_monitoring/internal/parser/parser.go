package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

func Parser(link string) (string, float32, error) {
	cmd := exec.Command("/usr/bin/python3", "./scraper.py", link)

	output, err := cmd.Output()
	if err != nil {
		return "", 0, err
	}

	var result struct {
		Name       string  `json:"name"`
		Sale_price float32 `json:"sale_price"`
		Error      string  `json:"error"`
	}

	err = json.Unmarshal(output, &result)
	if err != nil {
		return "", 0, fmt.Errorf("invalid json from parser: %w", err)
	}

	if result.Error != "" {
		if result.Error == "Товара нет в наличии" {
			return result.Name, 0, errors.New(result.Error)
		}
		return "", 0, fmt.Errorf("parser error: %s", result.Error)
	}

	return result.Name, result.Sale_price, nil
}
