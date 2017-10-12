package cmd

import (
	"fmt"
	"io/ioutil"
	"github.com/spf13/cobra"
	"github.com/mjolnir92/kdfs/restmsg"
	"github.com/vmihailenco/msgpack"
)

var storeCmd = &cobra.Command{
  Use:   "store",
  Short: "Store the file in the network",
  Long: `Stores the data in the given file in the network. The ID of the file is returned.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}
		req := restmsg.StoreRequest{File: content}
		// TODO: get host and port from some config
		url := "http://" + server + "/v1/store"
		b, err := postMsgPack(url, req)
		if err != nil {
			return err
		}
		var res restmsg.StoreResponse
		err = msgpack.Unmarshal(b, &res)
		if err != nil {
			return err
		}
		fmt.Println(res.ID)
		return nil
  },
}

func init() {
	RootCmd.AddCommand(storeCmd)
}
