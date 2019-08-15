package cmd

import (
	"heatbeat/pkg/aws"

	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:     "send",
	Short:   "send heartbeat SQS message",
	Example: "send",
	Run: func(cmd *cobra.Command, args []string) {
		aws.SendAliveSQSMessage()
	},
}
