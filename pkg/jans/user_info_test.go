package jans

import (
	"context"
	"testing"

	"github.com/moabu/kubelogin/pkg/kubeconfig"
)

func TestUserInfo(t *testing.T) {

	t.Skip("This test requires can only be run manually, as it requires a valid user token, retrieved via device authentication")

	encodedToken := ``

	token, err := kubeconfig.DecodeToken(encodedToken)
	if err != nil {
		t.Fatal(err)
	}

	accessToken := token.AccessToken

	client, err := NewClient(token.Url, token.ClientID, accessToken)
	if err != nil {
		t.Fatal(err)
	}

	ui, err := client.GetUserInfo(context.Background(), accessToken)
	if err != nil {
		t.Fatal(err)
	}

	if ui.UserName == "" {
		t.Error("expected user name, got empty string")
	}

}
