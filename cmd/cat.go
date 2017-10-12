package cmd

import (
	"encoding/binary"
	"os"
	"github.com/spf13/cobra"
	"github.com/vmihailenco/msgpack"
	"github.com/mjolnir92/kdfs/restmsg"
)

var catCmd = &cobra.Command{
  Use:   "cat",
  Short: "Read data with a specific ID and send to standard output",
  Long: `Read data with the given ID and send it to standard output. Unlike its namesake, it has nothing to do with concatenating files.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: get host and port from some config
		url := "http://" + server + "/v1/store/" + args[0]
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
	RootCmd.AddCommand(catCmd)
}
