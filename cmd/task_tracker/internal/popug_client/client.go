package popug_client

import (
	"fmt"
	"io"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"uber-popug/pkg/urls"
)

type client struct {
	endpoint string
	client   *http.Client
}

func New() *client {
	return &client{
		endpoint: urls.UsersUrl,
		client:   &http.Client{},
	}
}

func (c *client) GetAllPopugsIDs() ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/internal/popugs", c.endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
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
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		_, err = io.Copy(io.Discard, resp.Body)
		if err != nil {
			log.Printf("read response body: %s", err)
		}
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	var res []string
	err = jsoniter.ConfigFastest.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Printf("failed to decode response body: %s", err)
		return nil, fmt.Errorf("decode body: %w", err)
	}

	return res, nil
}
