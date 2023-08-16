package kubeconfig

import (
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func UpdateKubeconfig(path string, config *restclient.Config, token string) error {

	// read file
	apiCfg, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return err
	}

	// find current context
	currentContext := apiCfg.Contexts[apiCfg.CurrentContext]
	targetUser := currentContext.AuthInfo

	users := apiCfg.AuthInfos
	for userName, user := range users {
		if userName == targetUser {
			user.Token = token
			break
		}
	}

	clientcmd.WriteToFile(*apiCfg, path)

	return nil
}
