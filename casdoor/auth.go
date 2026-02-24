// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package casdoorsdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// AuthConfig is the core configuration.
// The first step to use this SDK is to use the InitConfig function to initialize the global authConfig.
type AuthConfig struct {
	Endpoint         string
	ClientId         string
	ClientSecret     string
	Certificate      string
	OrganizationName string
	ApplicationName  string
}

type Client struct {
	AuthConfig
	CustomHeaders map[string]string
}

// HttpClient interface has the method required to use a type as custom http client.
// The net/*http.Client type satisfies this interface.
type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Response struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Data2  interface{} `json:"data2"`
}

// client is a shared http Client.
var client HttpClient = &http.Client{}
var globalClient *Client

func InitConfig(endpoint string, clientId string, clientSecret string, certificate string, organizationName string, applicationName string) {
	globalClient = NewClient(endpoint, clientId, clientSecret, certificate, organizationName, applicationName)
}

func NewClient(endpoint string, clientId string, clientSecret string, certificate string, organizationName string, applicationName string) *Client {
	return NewClientWithConf(
		&AuthConfig{
			Endpoint:         endpoint,
			ClientId:         clientId,
			ClientSecret:     clientSecret,
			Certificate:      certificate,
			OrganizationName: organizationName,
			ApplicationName:  applicationName,
		})
}

func NewClientWithConf(config *AuthConfig) *Client {
	return &Client{
		AuthConfig:    *config,
		CustomHeaders: make(map[string]string),
	}
}

// SetHttpClient sets custom http Client.
func SetHttpClient(httpClient HttpClient) {
	client = httpClient
}

// OAuthOption is a function type for configuring OAuth requests.
type OAuthOption func(*oauthOptions)

// oauthOptions holds configuration options for OAuth operations.
type oauthOptions struct {
	httpClient *http.Client
}

// WithHTTPClient sets a custom http client for oauth operations.
func WithHTTPClient(httpClient *http.Client) OAuthOption {
	return func(opts *oauthOptions) {
		opts.httpClient = httpClient
	}
}

// GetOAuthToken gets the pivotal and necessary secret to interact with the Casdoor server
func (c *Client) GetOAuthToken(code string, state string, opts ...OAuthOption) (*oauth2.Token, error) {
	options := &oauthOptions{}
	for _, opt := range opts {
		opt(options)
	}

	config := oauth2.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("%s/api/login/oauth/authorize", c.Endpoint),
			TokenURL:  fmt.Sprintf("%s/api/login/oauth/access_token", c.Endpoint),
			AuthStyle: oauth2.AuthStyleInParams,
		},
		// RedirectURL: redirectUri,
		Scopes: nil,
	}

	ctx := context.Background()
	if options.httpClient != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, options.httpClient)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return token, err
	}

	if strings.HasPrefix(token.AccessToken, "error:") {
		return nil, errors.New(strings.TrimPrefix(token.AccessToken, "error: "))
	}

	return token, err
}

// RefreshOAuthToken refreshes the OAuth token
func (c *Client) RefreshOAuthToken(refreshToken string, opts ...OAuthOption) (*oauth2.Token, error) {
	options := &oauthOptions{}
	for _, opt := range opts {
		opt(options)
	}

	config := oauth2.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("%s/api/login/oauth/authorize", c.Endpoint),
			TokenURL:  fmt.Sprintf("%s/api/login/oauth/refresh_token", c.Endpoint),
			AuthStyle: oauth2.AuthStyleInParams,
		},
		// RedirectURL: redirectUri,
		Scopes: nil,
	}

	ctx := context.Background()
	if options.httpClient != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, options.httpClient)
	}

	token, err := config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken}).Token()
	if err != nil {
		return token, err
	}

	if strings.HasPrefix(token.AccessToken, "error:") {
		return nil, errors.New(strings.TrimPrefix(token.AccessToken, "error: "))
	}

	return token, err
}

// GetOAuthTokenWithClientCredentials gets access token using OAuth 2.0 client credentials grant flow
func (c *Client) GetOAuthTokenWithClientCredentials(scope string, opts ...OAuthOption) (*oauth2.Token, error) {
	options := &oauthOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Prepare form data
	formData := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     c.ClientId,
		"client_secret": c.ClientSecret,
	}
	if scope != "" {
		formData["scope"] = scope
	}

	// Create URL-encoded form
	contentType, body, err := createURLEncodedForm(formData)
	if err != nil {
		return nil, fmt.Errorf("failed to create form data: %w", err)
	}

	// Prepare request URL
	tokenURL := fmt.Sprintf("%s/api/login/oauth/access_token", c.Endpoint)

	// Use custom HTTP client if provided
	httpClient := client
	if options.httpClient != nil {
		httpClient = options.httpClient
	}

	// Create request
	req, err := http.NewRequest("POST", tokenURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	// Add custom headers
	for key, value := range c.CustomHeaders {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d, status: %s, body: %s", resp.StatusCode, resp.Status, string(respBytes))
	}

	// Parse token response
	var tokenResponse struct {
		AccessToken      string `json:"access_token"`
		TokenType        string `json:"token_type"`
		ExpiresIn        int64  `json:"expires_in"`
		Scope            string `json:"scope,omitempty"`
		RefreshToken     string `json:"refresh_token,omitempty"`
		Error            string `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
	}

	if err := json.Unmarshal(respBytes, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if tokenResponse.Error != "" {
		errorMsg := tokenResponse.Error
		if tokenResponse.ErrorDescription != "" {
			errorMsg = fmt.Sprintf("%s: %s", tokenResponse.Error, tokenResponse.ErrorDescription)
		}
		return nil, errors.New(errorMsg)
	}

	if strings.HasPrefix(tokenResponse.AccessToken, "error:") {
		return nil, errors.New(strings.TrimPrefix(tokenResponse.AccessToken, "error: "))
	}

	// Convert to oauth2.Token
	token := &oauth2.Token{
		AccessToken:  tokenResponse.AccessToken,
		TokenType:    tokenResponse.TokenType,
		RefreshToken: tokenResponse.RefreshToken,
	}

	// Set expiry time if provided
	if tokenResponse.ExpiresIn > 0 {
		token.Expiry = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	}

	return token, nil
}
