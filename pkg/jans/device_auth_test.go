package jans

import (
	"context"
	"os"
	"testing"
)

func TestDeviceAuth(t *testing.T) {

	host := os.Getenv("JANS_URL")
	user := os.Getenv("JANS_CLIENT_ID")
	pass := os.Getenv("JANS_CLIENT_SECRET")

	client, err := NewClient(host, user, pass)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.StartDeviceAuth(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if resp.DeviceCode == "" || resp.VerificationUriComplete == "" {
		t.Error("expected device code and url, got empty string(s)")
	}

	t.Logf("url: %s", resp.VerificationUriComplete)
	t.Logf("device code: %s", resp.DeviceCode)
	t.Fail()

}

func TestToken(t *testing.T) {

	t.Skip("This test requires can only be run manually, as it requires user interaction")

	host := os.Getenv("JANS_URL")
	user := os.Getenv("JANS_CLIENT_ID")
	pass := os.Getenv("JANS_CLIENT_SECRET")

	client, err := NewClient(host, user, pass)
	if err != nil {
		t.Fatal(err)
	}

	token, err := client.GetDeviceToken(context.Background(), "e3fee2fc933be19f841dcfaa208aaea986ad1f71a7cfe7d9")
	if err != nil {
		t.Fatal(err)
	}

	if token.AccessToken == "" {
		t.Error("expected access token, got empty string")
	}
}
