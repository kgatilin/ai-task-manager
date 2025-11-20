package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCompletionCommand creates a Cobra command for shell completion
func NewCompletionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion script",
		Long: `Generate shell completion script for tm.

To load completions:

Bash:
  $ source <(tm completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ tm completion bash > /etc/bash_completion.d/tm
  # macOS:
  $ tm completion bash > /usr/local/etc/bash_completion.d/tm

Zsh:
  $ source <(tm completion zsh)

  # To load completions for each session, add to ~/.zshrc:
  $ tm completion zsh > "${fpath[1]}/_tm"

Fish:
  $ tm completion fish | source

  # To load completions for each session:
  $ tm completion fish > ~/.config/fish/completions/tm.fish

PowerShell:
  PS> tm completion powershell | Out-String | Invoke-Expression

  # To load completions for each session, add to your PowerShell profile.`,
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			}
			return fmt.Errorf("unknown shell: %s", args[0])
		},
	}
}
