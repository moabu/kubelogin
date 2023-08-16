package cmd

import (
	"github.com/spf13/cobra"
	"github.com/moabu/kubelogin/pkg/jans"
)

// newProxyCmd provides a cobra command for running a Jans proxy
func newProxyCmd() *cobra.Command {

	var listenAddress *string
	var target *string

	cmd := &cobra.Command{
		Use:     "proxy --listen-address <listen-address> --target-url <target-url>",
		Example: `  $ kubelogin proxy --listen-address localhost:443 --target-url https://jans-proxy.jans.io`,
		Short:   "Start a Jans proxy",
		Long: `Start a Jans proxy

Runs a proxy that will forward requests from a Kubernetes cluster to a Jans 
server. The proxy will listen on the given address. Requests to the proxy
should be made using the Kubernetes webhook format.
		`,

		RunE: func(c *cobra.Command, args []string) error {

			if listenAddress == nil || target == nil {
				return c.Usage()
			}

			jans.StartHandler(*listenAddress, *target, opts.verbose)

			return nil
		},
	}

	listenAddress = cmd.Flags().StringP("listen-address", "l", "", "The network address to bind to.")
	cmd.MarkFlagRequired("listen-address")

	target = cmd.Flags().StringP("target-url", "t", "", "The URL of the target Jans instance.")
	cmd.MarkFlagRequired("target-url")

	return cmd
}
