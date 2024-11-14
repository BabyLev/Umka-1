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
	SatelliteId int       `json:"satelliteId"`
	Date        time.Time `json:"date"`
	Name        string    `json:"name"`
	Line1       string    `json:"line1"`
	Line2       string    `json:"line2"`
}

func New(url string) *Client {
	return &Client{
		url: url,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) GetSatelliteInfo(ctx context.Context, noradID int) (R4uabSat, error) {
	reqURL := fmt.Sprintf("%s/tle/%d", c.url, noradID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return R4uabSat{}, fmt.Errorf("ошибка составления запроса: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return R4uabSat{}, fmt.Errorf("request error: %w", err)
	}

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
