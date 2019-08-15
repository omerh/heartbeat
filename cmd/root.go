package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCommandExample = `
	# Run hearbeat send
	heartbeat send
	# Run heartbeat check
	heartbeat check
`
)

var rootCmd = &cobra.Command{
	Use:     "heartbeat",
	Short:   "Heartbeat for send heartbeat message to AWS SQS and Check if alive",
	Example: rootCommandExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(checkCmd)
}

// Execute usign cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
