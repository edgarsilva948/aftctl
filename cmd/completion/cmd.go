/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package completion

import (
	"os"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates completion scripts",
	Long: `To load completions:

Bash:

  $ source <(aftctl completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ aftctl completion bash > /etc/bash_completion.d/aftctl
  # macOS:
  $ aftctl completion bash > /usr/local/etc/bash_completion.d/aftctl

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ aftctl completion zsh > "${fpath[1]}/_aftctl"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ aftctl completion fish | source

  # To load completions for each session, execute once:
  $ aftctl completion fish > ~/.config/fish/completions/aftctl.fish

PowerShell:

  PS> aftctl completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> aftctl completion powershell > aftctl.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Default to bash for backwards compatibility
			cmd.Root().GenBashCompletion(os.Stdout)
			return
		}
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}
