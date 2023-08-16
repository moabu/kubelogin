package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate commnand line completion script",
		Long: `To load completions:
	
Bash:

$ source <(kubelogin completion bash)

# To load completions for each session, execute once:
Linux:
	$ kubelogin completion bash > /etc/bash_completion.d/kubelogin
MacOS:
	$ kubelogin completion bash > /usr/local/etc/bash_completion.d/kubelogin

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ kubelogin completion zsh > "${fpath[1]}/_kubelogin"

# You will need to start a new shell for this setup to take effect.

Fish:

$ kubelogin completion fish | source

# To load completions for each session, execute once:
$ kubelogin completion fish > ~/.config/fish/completions/kubelogin.fish
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	}
}
