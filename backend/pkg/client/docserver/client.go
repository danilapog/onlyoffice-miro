package docserver

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/pkg/common"
)

type client struct {
	httpClient *http.Client
}

func NewClient() Client {
	return &client{
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
			Transport: common.NewRetryableTransport(&http.Transport{
				MaxIdleConnsPerHost:    100,
				IdleConnTimeout:        90 * time.Second,
				MaxResponseHeaderBytes: 1 << 20,
				DisableCompression:     false,
				ForceAttemptHTTP2:      true,
			},
			),
		},
	}
}

func generateRandomKey(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}

func (c *client) GetServerVersion(ctx context.Context, base string, opts ...Option) (*ServerVersionResponse, error) {
	options := DefaultClientOptions()
	for _, opt := range opts {
		opt(options)
	}

	address := strings.TrimRight(base, "/")
	_, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("malformed docserver address: %w", err)
	}

	body := GetServerVersionRequest{C: "version", Token: options.Token}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal version request body: %w", err)
	}

	url := common.Concat(address, "/command", "?shardKey=", generateRandomKey(8))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if options.Header != "" {
		req.Header.Set(options.Header, options.Token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response ServerVersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
