package jans

import (
	"context"
	"fmt"
	"log"
	"net/url"
)

type DeviceAuthResponse struct {
	UserCode                string `json:"user_code"`
	DeviceCode              string `json:"device_code"`
	Interval                int    `json:"interval"`
	VerificationUri         string `json:"verification_uri"`
	VerificationUriComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *Client) StartDeviceAuth(ctx context.Context) (*DeviceAuthResponse, error) {

	req := url.Values{}
	req.Add("client_id", c.clientId)
	req.Add("scope", `openid+profile+username+email+offline_access`)

	var resp DeviceAuthResponse

	params := requestParams{
		method:      "POST",
		path:        "/jans-auth/restv1/device_authorization",
		accept:      "application/json",
		contentType: "application/x-www-form-urlencoded",
		payload:     []byte(req.Encode()),
		resp:        &resp,
	}

	if err := c.request(ctx, params); err != nil {
		werr := fmt.Errorf("error sending device auth request to server: %w", err)
		log.Println(werr)
		return nil, werr
	}

	return &resp, nil
}

func (c *Client) GetDeviceToken(ctx context.Context, deviceCode string) (*TokenResponse, error) {

	req := url.Values{}
	req.Add("client_id", c.clientId)
	req.Add("scope", `openid+profile+username+email+offline_access`)
	req.Add("grant_type", `urn:ietf:params:oauth:grant-type:device_code`)
	req.Add("grant_type", `refresh_token`)
	req.Add("device_code", deviceCode)

	var ret TokenResponse

	params := requestParams{
		method:      "POST",
		path:        "/jans-auth/restv1/token",
		accept:      "application/json",
		contentType: "application/x-www-form-urlencoded",
		payload:     []byte(req.Encode()),
		resp:        &ret,
	}

	if err := c.request(ctx, params); err != nil {
		return nil, fmt.Errorf("error sending device auth request to server: %w", err)
	}

	return &ret, nil
}
