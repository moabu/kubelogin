package kubeconfig

import "testing"

func TestToken(t *testing.T) {

	token := TokenData{
		Url:            "https://jans-instance.io",
		ClientID:       "1900.d4a64508-b347-4dc0-beb7-d85a737d8784",
		ClientPassword: "xxxxxxx",
		AccessToken:    "fake_token",
	}

	encoded, err := EncodeToken(token)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := DecodeToken(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.Url != token.Url {
		t.Errorf("incorrect URL, expected %s, got %s", token.Url, decoded.Url)
	}

	if decoded.ClientID != token.ClientID {
		t.Errorf("incorrect client_id, expected %s, got %s", token.ClientID, decoded.ClientID)
	}

	if decoded.ClientPassword != token.ClientPassword {
		t.Errorf("incorrect client_password, expected %s, got %s", token.ClientPassword, decoded.ClientPassword)
	}

	if decoded.AccessToken != token.AccessToken {
		t.Errorf("incorrect access_token, expected %s, got %s", token.AccessToken, decoded.AccessToken)
	}
}
