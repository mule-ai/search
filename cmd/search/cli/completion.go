package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// AddCompletionCommand adds the shell completion command to the root command.
//
// The completion command generates shell-specific completion scripts for
// bash, zsh, and fish.
func AddCompletionCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "completion [bash|zsh|fish]",
		Short: "Generate shell completion scripts",
		Long: `To load completions in your current shell session:

  Bash:
    source <(search completion bash)

  Zsh:
    source <(search completion zsh)

  Fish:
    search completion fish | source

To install completions system-wide:

  Bash (Linux):
    search completion bash > /etc/bash_completion.d/search

  Zsh:
    search completion zsh > ~/.zfunc/_search

  Fish:
    search completion fish > ~/.config/fish/completions/search.fish`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				return rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				return rootCmd.GenFishCompletion(os.Stdout, true)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	})
}
