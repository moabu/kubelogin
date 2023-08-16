package cmd

import (
	"github.com/spf13/cobra"
)

type rootOptions struct {
	verbose bool
}

var opts rootOptions

func NewRootCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:          "kubelogin",
		Short:        "kubelogin is a Kubernetes credentials plugin implementing OpenID Connect (OIDC) authentication",
		SilenceUsage: true,
	}

	opts = rootOptions{}

	cmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "Enable verbose output.")

	cmd.AddCommand(newGetTokenCmd())
	cmd.AddCommand(newCreateTokenCmd())
	cmd.AddCommand(newProxyCmd())
	cmd.AddCommand(newCompletionCmd())

	return cmd
}
