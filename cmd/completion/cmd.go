/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package completion

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Cmd represents the completion command for generating shell completion scripts.
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
	Run:                   Run,
}

// Run is the main function for the 'completion' command. It delegates to RunCompletion.
func Run(cmd *cobra.Command, args []string) {
	RunCompletion(cmd, args, os.Stdout)
}

// RunCompletion generates the shell completion script based on the provided arguments.
func RunCompletion(cmd *cobra.Command, args []string, out io.Writer) {
	if len(args) == 0 {
		// Default to bash for backwards compatibility
		generateBashCompletion(cmd.Root(), out)
		return
	}
	switch args[0] {
	case "bash":
		generateBashCompletion(cmd.Root(), out)
	case "zsh":
		generateZshCompletion(cmd.Root(), out)
	case "powershell":
		generatePowerShellCompletion(cmd.Root(), out)
	}
}

// generateBashCompletion generates Bash completion script.
func generateBashCompletion(cmd *cobra.Command, out io.Writer) {
	cmd.GenBashCompletion(out)
}

// generateZshCompletion generates Zsh completion script.
func generateZshCompletion(cmd *cobra.Command, out io.Writer) {
	cmd.GenZshCompletion(out)
}

// generatePowerShellCompletion generates PowerShell completion script.
func generatePowerShellCompletion(cmd *cobra.Command, out io.Writer) {
	cmd.GenPowerShellCompletionWithDesc(out)
}
