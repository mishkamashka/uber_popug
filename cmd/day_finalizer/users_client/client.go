package accounting_client

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type client struct {
	endpoint string
	client   *http.Client
}

func New() *client {
	return &client{
		endpoint: "http://localhost:2400/api",
		client:   &http.Client{},
	}
}

type CheckoutRequest struct{}

func (c *client) GetUserEmail(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("empty userID")
	}

	url := fmt.Sprintf("%s/internal/checkout?user_id=%s", c.endpoint, userID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("close response body: %s", closeErr)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		_, err = io.Copy(io.Discard, resp.Body)
		if err != nil {
			log.Printf("read response body: %s", err)
		}
		return "", errors.New("user not found")
	}

	if resp.StatusCode != http.StatusOK {
		_, err = io.Copy(io.Discard, resp.Body)
		if err != nil {
			log.Printf("read response body: %s", err)
		}
		return "", fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	res := struct {
		Email string `json:"email"`
	}{}

	err = jsoniter.ConfigFastest.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Printf("failed to decode response body: %s", err)
		return "", fmt.Errorf("decode body: %w", err)
	}

	return res.Email, nil
}
