package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func Run(rootCmd *cobra.Command) error {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		// EOF (Ctrl+D)
		if !scanner.Scan() {
			return nil
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if line == "q" || line == "quit" {
			return nil
		}

		args := strings.Fields(line)

		rootCmd.SetArgs(args)

		if err := rootCmd.Execute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
