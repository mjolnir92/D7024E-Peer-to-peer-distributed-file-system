package cmd

import (
	"fmt"
	"io/ioutil"
	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
  Use:   "pin",
  Short: "Protect an ID from deletion",
  Long: `The pin command makes sure important data is not deleted.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("pinCmd.Run!")
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
	RootCmd.AddCommand(pinCmd)
}
