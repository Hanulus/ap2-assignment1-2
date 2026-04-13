package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// PaymentHTTPClient calls the Payment Service over REST
type PaymentHTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewPaymentHTTPClient(baseURL string) *PaymentHTTPClient {
	return &PaymentHTTPClient{
		baseURL: baseURL,
		// Timeout prevents hanging when Payment Service is down
		client: &http.Client{Timeout: 2 * time.Second},
	}
}

type authorizeRequest struct {
	OrderID string `json:"order_id"`
	Amount  int64  `json:"amount"`
}

type authorizeResponse struct {
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
}

// Authorize calls POST /payments on the Payment Service
func (c *PaymentHTTPClient) Authorize(orderID string, amount int64) (string, error) {
	body, _ := json.Marshal(authorizeRequest{OrderID: orderID, Amount: amount})

	resp, err := c.client.Post(c.baseURL+"/payments", "application/json", bytes.NewReader(body))
	if err != nil {
		// Network error or timeout — return 503-friendly error
		return "", errors.New("payment service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("payment declined (status %d)", resp.StatusCode)
	}

	var result authorizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Status != "Authorized" {
		return "", errors.New("payment declined")
	}

	return result.TransactionID, nil
}
