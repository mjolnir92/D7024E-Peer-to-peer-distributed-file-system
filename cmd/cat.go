
package cmd

import (
	"fmt"
	"io/ioutil"
	"binary"
	"os"
	"github.com/spf13/cobra"
)

var catCmd = &cobra.Command{
  Use:   "cat",
  Short: "Read data with a specific ID and send to standard output",
  Long: `Read data with the given ID and send it to standard output. Unlike its namesake, it has nothing to do with concatenating files.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("catCmd.Run!")
		// TODO: get host and port from some config
		url := "http://localhost:8080" + "/v1/store/" + args[0]
		b, err := get(url)
		if err != nil {
			return err
		}
		var res restmsg.CatResponse
		err = msgpack.Unmarshal(b, &res)
		if err != nil {
			return err
		}
		err = binary.Write(os.Stdout, binary.LittleEndian, res.File)
		if err != nil {
			return err
		}
		return err
  },
}

func init() {
	RootCmd.AddCommand(storeCmd)
}




var storeCmd = &cobra.Command{
  Use:   "store",
  Short: "Store the file in the network",
  Long: `Stores the data in the given file in the network. The ID of the file is returned.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("storeCmd.Run!")
		content, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}
		req := restmsg.StoreRequest{File: content}
		// TODO: get host and port from some config
		url := "http://localhost:8080" + "/v1/store"
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
