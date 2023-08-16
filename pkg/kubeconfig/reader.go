package kubeconfig

import (
	"fmt"
	"os"

	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func ReadKubeconfig(kubeconfigPath string) (*restclient.Config, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("could not build config from kubeconfig location '%s': %w", kubeconfigPath, err)
	}

	return config, nil
}

func GetKubeconfigLocation(path string) (string, error) {

	// first try to get the right location
	location := ""
	if path != "" {
		location = path
	}

	// fallback to environment variable
	if location == "" {
		location = os.Getenv("KUBECONFIG")
	}

	// fallback to default location
	if location == "" {
		location = os.Getenv("HOME") + "/.kube/config"
	}

	if location == "" {
		return "", fmt.Errorf("no valid kubeconfig location was found or provided")
	}

	// check if file exists and is readable
	if _, err := os.Stat(location); os.IsNotExist(err) {
		return "", fmt.Errorf("kubeconfig file does not exist: %s", location)
	}

	return location, nil
}
