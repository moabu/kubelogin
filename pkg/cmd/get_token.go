package cmd

import (
	"bufio"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/moabu/kubelogin/pkg/jans"
	"github.com/moabu/kubelogin/pkg/kubeconfig"
)

// newGetTokenCmd provides a cobra command for starting a device authorization grant
func newGetTokenCmd() *cobra.Command {

	var kubeconfigPath *string

	cmd := &cobra.Command{
		Use:                   "get-token --kubeconfig <kubeconfig>",
		Example:               `  $ kubelogin get-token --kubeconfig ~/.kube/config`,
		Short:                 "Get a token via device authorization grant",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE: func(c *cobra.Command, args []string) error {

			path := ""
			if kubeconfigPath != nil {
				path = *kubeconfigPath
			}

			location, err := kubeconfig.GetKubeconfigLocation(path)
			if err != nil {
				if opts.verbose {
					return fmt.Errorf("could not get kubeconfig location: %v", err)
				}
				return fmt.Errorf("could not find kubeconfig location, or no kubeconfig file was provided")
			}

			config, err := kubeconfig.ReadKubeconfig(path)
			if err != nil {
				if opts.verbose {
					return fmt.Errorf("could not read kubeconfig file from location %s: %v", location, err)
				}
				return fmt.Errorf("could not read kubeconfig file")
			}

			if opts.verbose {
				fmt.Printf("Successfully read kubeconfig from location '%s'\n", location)
			}

			originalToken, err := kubeconfig.DecodeToken(config.BearerToken)
			if err != nil {
				if opts.verbose {
					return fmt.Errorf("could not parse token in kubeconfig file: %v", err)
				}
				return fmt.Errorf("could not parse token in kubeconfig file")
			}

			if opts.verbose {
				fmt.Printf("Using client %s for jans instance %s\n", originalToken.ClientID, originalToken.Url)
			}
			cl, err := jans.NewClient(originalToken.Url, originalToken.ClientID, originalToken.ClientPassword)
			if err != nil {
				return err
			}

			// trigger the device authorization grant
			resp, err := cl.StartDeviceAuth(c.Context())
			if err != nil {
				return fmt.Errorf("could not start device authorization grant: %w", err)
			}

			c.Printf("Please visit '%s' to authenticate\n", resp.VerificationUriComplete)
			c.Printf("Once done, press enter to continue\n")

			scanner := bufio.NewScanner(c.InOrStdin())
			for scanner.Scan() {
				break
			}

			// now get token from Jans, using the device code
			tokenResp, err := cl.GetDeviceToken(c.Context(), resp.DeviceCode)
			if err != nil {
				return fmt.Errorf("could not get token: %w", err)
			}

			if opts.verbose {
				fmt.Printf("Successfully got access token from Jans: %s\n", tokenResp.AccessToken)
			}

			// generate new token
			newToken := kubeconfig.TokenData{
				Url:            originalToken.Url,
				ClientID:       originalToken.ClientID,
				ClientPassword: originalToken.ClientPassword,
				AccessToken:    tokenResp.AccessToken,
			}

			newTokenEncoded, err := kubeconfig.EncodeToken(newToken)
			if err != nil {
				return fmt.Errorf("could not encode token: %w", err)
			}

			// update kubeconfig
			if err := kubeconfig.UpdateKubeconfig(location, config, newTokenEncoded); err != nil {
				if opts.verbose {
					return fmt.Errorf("could not update kubeconfig file at location %s: %v", location, err)
				}
				return fmt.Errorf("could not update kubeconfig file")
			}

			if opts.verbose {
				fmt.Printf("Successfully updated kubeconfig file at location '%s' with new access token\n", location)
			}

			return nil
		},
	}

	kubeconfigPath = cmd.Flags().StringP("kubeconfig", "k", "", "Path to kubeconfig file. If not provided, the KUBECONFIG environment variable will be used or fallback to '~/.kube/config'.")
	cmd.MarkFlagFilename("kubeconfig")

	return cmd
}
