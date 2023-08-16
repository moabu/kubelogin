package kubeconfig

import (
	"os"
	"testing"
)

func TestReader(t *testing.T) {

	file, err := createTempKubeconfig()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	client, err := ReadKubeconfig(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if client.Host != "https://jans-k8s-instance.io:6443" {
		t.Errorf("expected %s, got %s", "https://jans-k8s-instance.io:6443", client.Host)
	}
}

func TestKubeconfigLocation(t *testing.T) {

	file, err := createTempKubeconfig()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	location, err := GetKubeconfigLocation(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if location != file.Name() {
		t.Errorf("expected %s, got %s", file.Name(), location)
	}

	envFile, err := createTempKubeconfig()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(envFile.Name())

	os.Setenv("KUBECONFIG", envFile.Name())

	envLocation, err := GetKubeconfigLocation(envFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if envLocation != envFile.Name() {
		t.Errorf("expected %s, got %s", file.Name(), envLocation)
	}

}

func createTempKubeconfig() (*os.File, error) {

	kubeconfigData := `apiVersion: v1
	clusters:
	- cluster:
			server: https://jans-k8s-instance.io:6443
		name: test-cluster
	contexts:
	- context:
			cluster: test-cluster
			user: johndoe
		name: johndoe@test-cluster
	current-context: johndoe@test-cluster
	kind: Config
	preferences: {}
	users:
	- name: johndoe
		user:
			token: eyJ1cmwiOiJodHRwczovL2phbnMtaW5zdGFuY2UuaW8iLCJjbGllbnQiOiIxOTAwLmQ0YTY0NTA4LWIzNDctNGRjMC1iZWI3LWQ4NWE3MzdkODc4NCIsInBhc3N3b3JkIjoieHh4eHh4eCJ9
	`
	// write data to temp file
	file, err := os.CreateTemp("", "kubeconfig-*.yaml")
	if err != nil {
		return nil, err
	}

	if _, err := file.WriteString(kubeconfigData); err != nil {
		return nil, err
	}

	return file, nil
}
