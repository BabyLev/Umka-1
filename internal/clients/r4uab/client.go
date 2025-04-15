package r4uab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TODO: написать http клиент, который сделает запрос на апи и вернет информацию о спутнике

type Client struct {
	url    string // "https://api.r4uab.ru"
	client *http.Client
}

type R4uabSat struct {
	SatelliteId string `json:"satelliteId"`
	Date        string `json:"date"`
	Name        string `json:"name"`
	Line1       string `json:"line1"`
	Line2       string `json:"line2"`
}

func New(url string) *Client {
	return &Client{
		url: url,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) GetSatelliteInfo(ctx context.Context, noradID int64) (R4uabSat, error) {
	if noradID <= 0 {
		return R4uabSat{}, fmt.Errorf("noradID должен быть больше 0")
	}

	reqURL := fmt.Sprintf("%s/tle/%d", c.url, noradID) // https://api.r4uab.ru/tle/12345/

	// формируем HTTP запрос, используя: метод GET и адрес, тело запроса пустое
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return R4uabSat{}, fmt.Errorf("ошибка составления запроса: %w", err)
	}

	// устанавливает соединение с сервером, отправляет на сервер HTTP запрос, ждет обротку, получает ответ и кладет в переменную resp
	resp, err := c.client.Do(req)
	if err != nil {
		return R4uabSat{}, fmt.Errorf("request error: %w", err)
	}

	// закрываем соединение после прочтения ответа от сервера(59 строка)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return R4uabSat{}, fmt.Errorf("expected status 200, but got: %d", resp.StatusCode)
	}

	var sat R4uabSat

	err = json.NewDecoder(resp.Body).Decode(&sat)
	if err != nil {
		return R4uabSat{}, err
	}

	return sat, nil
}
