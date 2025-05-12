package oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/config"
)

type client[T any] struct {
	config     *config.OAuthConfig
	httpClient *http.Client
	errors     *Errors
}

func NewOAuthClient[T any](config *config.OAuthConfig) (OAuthClient[T], error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &client[T]{
		config: config,
		errors: NewErrors(),
		httpClient: &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				MaxIdleConnsPerHost:    50,
				IdleConnTimeout:        90 * time.Second,
				ResponseHeaderTimeout:  1500 * time.Millisecond,
				MaxResponseHeaderBytes: 1 << 20,
				DisableCompression:     false,
				ForceAttemptHTTP2:      true,
			},
		},
	}, nil
}

func (c *client[T]) buildFormData(params map[string]string) url.Values {
	data := url.Values{}
	for key, value := range params {
		data.Set(key, value)
	}

	return data
}

func (c *client[T]) doRequest(ctx context.Context, url string, data url.Values) (T, error) {
	var response T

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return response, c.errors.FailedToCreateRequest(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return response, c.errors.FailedToSendRequest(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return response, c.errors.RequestFailed(resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, c.errors.FailedToDecodeResponse(err)
	}

	return response, nil
}

func (c *client[T]) Exchange(ctx context.Context, code string) (T, error) {
	var zero T
	data := c.buildFormData(map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     c.config.ClientID,
		"client_secret": c.config.ClientSecret,
		"redirect_uri":  c.config.RedirectURI,
	})

	address, err := url.Parse(c.config.TokenURI)
	if err != nil {
		return zero, c.errors.FailedToExchangeToken(err)
	}

	query := address.Query()
	query.Set("client_id", c.config.ClientID)
	query.Set("code", code)
	query.Set("redirect_uri", c.config.RedirectURI)
	address.RawQuery = query.Encode()

	response, err := c.doRequest(ctx, address.String(), data)
	if err != nil {
		return zero, c.errors.FailedToExchangeToken(err)
	}

	return response, nil
}

func (c *client[T]) Refresh(ctx context.Context, refreshToken string) (T, error) {
	var zero T
	data := c.buildFormData(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     c.config.ClientID,
		"client_secret": c.config.ClientSecret,
	})

	address, err := url.Parse(c.config.TokenURI)
	if err != nil {
		return zero, c.errors.FailedToRefreshToken(err)
	}

	query := address.Query()
	query.Set("grant_type", "refresh_token")
	query.Set("refresh_token", refreshToken)
	address.RawQuery = query.Encode()

	response, err := c.doRequest(ctx, address.String(), data)
	if err != nil {
		return zero, c.errors.FailedToRefreshToken(err)
	}

	return response, nil
}
