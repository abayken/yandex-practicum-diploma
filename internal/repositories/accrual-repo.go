package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AccrualRepository struct {
	BaseURL string
}

type OrderInfo struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (repo AccrualRepository) FetchOrderInfo(orderNumber string) (*OrderInfo, error) {
	//return &OrderInfo{Number: "6387478398972289", Status: "PROCESSED", Accrual: 500}, nil
	url := fmt.Sprintf("%s/api/orders/%s", repo.BaseURL, orderNumber)
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var order OrderInfo

	err = json.NewDecoder(response.Body).Decode(&order)

	if err != nil {
		return nil, err
	}

	return &order, nil
}
