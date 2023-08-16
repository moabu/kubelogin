package kubeconfig

import (
	"os"
	"testing"
)

func TestUpdater(t *testing.T) {

	file, err := createTempKubeconfig()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	token := TokenData{
		Url:         "https://jans-instance.io",
		AccessToken: "foobar",
	}
	encodedToken, err := EncodeToken(token)
	if err != nil {
		t.Fatalf("failed to encode token: %v", err)
	}

	client, err := ReadKubeconfig(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	err = UpdateKubeconfig("/Users/vhristov/Work/upwork/tf-provider/kubelogin/kubeconfig.yaml", client, encodedToken)
	if err != nil {
		t.Fatal(err)
	}

	// read client again
	client, err = ReadKubeconfig("/Users/vhristov/Work/upwork/tf-provider/kubelogin/test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if client.BearerToken != encodedToken {
		t.Fatalf("expected token %s, got %s", encodedToken, client.BearerToken)
	}

}
