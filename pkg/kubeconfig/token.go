package kubeconfig

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type TokenData struct {
	Url            string `json:"url,omitempty"`
	ClientID       string `json:"client_id,omitempty"`
	ClientPassword string `json:"client_password,omitempty"`
	AccessToken    string `json:"access_token,omitempty"`
}

// DecodeToken parses the provided bearer token and returns the token data
// that can be used to authenticate a user via Jans.
func DecodeToken(token string) (*TokenData, error) {

	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	// base64 decode token
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("error decoding token: %w", err)
	}

	var tokenData TokenData
	if err := json.Unmarshal(decoded, &tokenData); err != nil {
		return nil, fmt.Errorf("error unmarshalling token: %w", err)
	}

	return &tokenData, nil
}

// EncodeToken encodes the provided token data and returns a base64 encoded
// bearer token that can be used in a kubeconfig file to authenticate a user
// via Jans.
func EncodeToken(token TokenData) (string, error) {

	b, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("error marshalling token: %w", err)
	}

	// base64 encode token
	encoded := base64.StdEncoding.EncodeToString(b)

	return encoded, nil
}
