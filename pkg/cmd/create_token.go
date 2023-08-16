package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/moabu/kubelogin/pkg/kubeconfig"
)

// newCreateTokenCmd provides a cobra command for creating a valid bearer token
// that can be used in a kubeconfig file to authenticate a user via Jans.
func newCreateTokenCmd() *cobra.Command {

	var (
		url      *string
		client   *string
		password *string
	)

	cmd := &cobra.Command{
		Use:                   "create-token --url <jans-url> --client <client-id> --password <client-password>",
		Example:               `  $ kubelogin create-token --url https://jans-instance.io --client 1900.d4a64508-b347-4dc0-beb7-d85a737d8784 --password xxxxxxx`,
		Short:                 "Create a token for the provided Jans credentials",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE: func(c *cobra.Command, args []string) error {

			if url == nil || client == nil || password == nil {
				return c.Usage()
			}

			token := kubeconfig.TokenData{
				Url:            *url,
				ClientID:       *client,
				ClientPassword: *password,
			}

			encoded, err := kubeconfig.EncodeToken(token)
			if err != nil {
				return fmt.Errorf("error encoding token: %w", err)
			}

			c.OutOrStdout().Write([]byte(encoded))
			c.OutOrStdout().Write([]byte("\n"))

			return nil
		},
	}

	url = cmd.Flags().StringP("url", "u", "", "URL of the jans instance against which the authentication should be performed")
	cmd.MarkFlagRequired("url")

	client = cmd.Flags().StringP("client", "c", "", "Client ID of the client that should be used for authentication")
	cmd.MarkFlagRequired("client")

	password = cmd.Flags().StringP("password", "p", "", "Password of the client that should be used for authentication")
	cmd.MarkFlagRequired("password")

	return cmd
}
