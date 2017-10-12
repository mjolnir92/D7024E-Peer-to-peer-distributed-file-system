package cmd

import (
	"github.com/spf13/cobra"
)

var unpinCmd = &cobra.Command{
  Use:   "unpin",
  Short: "Remove the pin status of an ID",
  Long: `Unpin allows the data to be deleted again.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: get host and port from some config
		url := "http://" + server + "/v1/unpin/" + args[0]
		_, err := postNoBody(url)
		if err != nil {
			return err
		}
		return nil
  },
}

func init() {
	RootCmd.AddCommand(unpinCmd)
}
