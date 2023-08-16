# kubelogin

A [Kubernetes client-go credential plugin](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#client-go-credential-plugins) implementing OIDC authentication.

The tool can be used to trigger a device flow authentication request and store the resulting token in the current kubeconfig file.

Available commads are:

|Command|Description|
|-------|-----------|
|`completion`|Generate completion script for the specified shell.|
|`create-token`|Creates a token that can be embedded in a kubeconfig file. This would usually be used on the server, before providing the kubeconfig file to a user.|
|`get-token`|Trigger a device flow authentication request and store the resulting token in the current kubeconfig file.|
|`proxy`|Start a proxy server that can be used to intercept Kubernetes webhook authentication calls and forward them to Jans and respectively map the response to the format expected by Kubernetes.|
