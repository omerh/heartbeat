package cmd

import (
	"heartbeat/pkg/aws"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:     "check",
	Short:   "check heartbeat SQS message",
	Example: "check",
	Run: func(cmd *cobra.Command, args []string) {
		aws.CheckAliveSQSMessage()
	},
}
