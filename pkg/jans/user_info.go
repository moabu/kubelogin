package jans

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
)

type UserInfoResponse struct {
	UserName string `json:"name"`
	UID      string `json:"inum"`
	Email    string `json:"email"`
}

func (c *Client) GetUserInfo(ctx context.Context, token string) (*UserInfoResponse, error) {

	req := url.Values{}
	req.Add("access_token", token)

	var resp string

	params := requestParams{
		method:      "POST",
		path:        "/jans-auth/restv1/userinfo",
		accept:      "application/json",
		contentType: "application/x-www-form-urlencoded",
		token:       token,
		payload:     []byte(req.Encode()),
		resp:        &resp,
		returnRaw:   true,
	}

	if err := c.request(ctx, params); err != nil {
		werr := fmt.Errorf("error sending device auth request to server: %w", err)
		log.Println(werr)
		return nil, werr
	}

	// decode JWT token
	// var claims map[string]any
	jwtPayload, err := parseJWT(resp)
	if err != nil {
		return nil, err
	}

	var ret UserInfoResponse

	if err := json.Unmarshal(jwtPayload, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func parseJWT(p string) ([]byte, error) {
	parts := strings.Split(p, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("malformed jwt, expected 3 parts got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("malformed jwt payload: %v", err)
	}
	return payload, nil
}
