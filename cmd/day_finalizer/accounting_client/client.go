package accounting_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"uber-popug/pkg/accounting"
	"uber-popug/pkg/urls"
)

type client struct {
	endpoint string
	client   *http.Client
}

func New() *client {
	return &client{
		endpoint: urls.AccountingUrl,
		client:   &http.Client{},
	}
}

type CheckoutRequest struct{}

func (c *client) Checkout(userID string, dayTotal int) error {
	reqData := accounting.CheckoutRequest{
		UserID:   userID,
		DayTotal: dayTotal,
	}

	body, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/internal/checkout", c.endpoint), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request: %w", err)
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
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		_, err = io.Copy(io.Discard, resp.Body)
		if err != nil {
			log.Printf("read response body: %s", err)
		}
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	return nil
}
