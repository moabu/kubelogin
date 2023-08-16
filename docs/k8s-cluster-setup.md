# Setting up a Kubernetes Cluster with Kubeadm

To be able to test the kubelogin plugin, you need a Kubernetes cluster. This guide will help you set up a Kubernetes cluster using Kubeadm.

## Installation

Prepare instance following [this guide](https://blog.radwell.codes/2022/07/single-node-kubernetes-cluster-via-kubeadm-on-ubuntu-22-04/), stopping before the "Create the cluster using kubeadm" section.

Create the following files:

`/etc/authn-config.yaml`:

```yaml
apiVersion: v1
kind: Config
clusters:
  - name: authn
    cluster:
      server: <JANS-PROXY-DNS-NAME>
      insecure-skip-tls-verify: true
users:
  - name: kube-apiserver
contexts:
- context:
    cluster: authn
    user: kube-apiserver
  name: authn
current-context: authn
```

`kubeadm-config.yaml`:

```yaml
apiVersion: kubeadm.k8s.io/v1beta3
bootstrapTokens:
- groups:
  - system:bootstrappers:kubeadm:default-node-token
  token: abcdef.0123456789abcdef
  ttl: 24h0m0s
  usages:
  - signing
  - authentication
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: 1.2.3.4
  bindPort: 6443
nodeRegistration:
  criSocket: unix:///var/run/containerd/containerd.sock
  imagePullPolicy: IfNotPresent
  name: node
  taints: null
---
apiServer:
  certSANs:
  - <CLUSTER-DNS-NAME>
  extraArgs:
    authentication-token-webhook-config-file: /etc/authn-config.yaml
  extraVolumes:
  - hostPath: /root/authn-config.yaml
    mountPath: /etc/authn-config.yaml
    name: authentication-token-webhook-config-file
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta3
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns: {}
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: registry.k8s.io
kind: ClusterConfiguration
kubernetesVersion: 1.27.0
networking:
  dnsDomain: cluster.local
  serviceSubnet: 10.96.0.0/12
scheduler: {}

```

Then configure the cluster using the following command:

```bash
kubeadm init --config kubeadm-config.yaml
```

Untaint the master node:

```bash
kubectl taint nodes --all node-role.kubernetes.io/control-plane-
```

Install networking following the instructions on [Cilium](https://docs.cilium.io/en/stable/gettingstarted/k8s-install-default/)
