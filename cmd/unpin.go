
package cmd

import (
	"fmt"
	"io/ioutil"
	"github.com/spf13/cobra"
)

var unpinCmd = &cobra.Command{
  Use:   "unpin",
  Short: "Remove the pin status of an ID",
  Long: `Unpin allows the data to be deleted again.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("unpinCmd.Run!")
		// TODO: get host and port from some config
		url := "http://localhost:8080" + "/v1/pin/" + args[0]
		b, err := postMsgPack(url)
		if err != nil {
			return err
		}
		return nil
  },
}

func init() {
	RootCmd.AddCommand(unpinCmd)
}
