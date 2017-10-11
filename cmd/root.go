package cmd

import (
	"github.com/spf13/cobra"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
	"os"
)

var RootCmd = &cobra.Command{
  Use:   "kdfs",
  Short: "kdfs is a DHT-based distributed file system",
  Long: `Interfaces with a kdfs server to store and retrieve files from the DFS.`,
  Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
  },
}

func init() {}
func Execute() {
	RootCmd.Execute()
}

// postMsgPack sends a post request with a msgpack-encoded body. The response body is returned.
func postMsgPack(url string, body []byte) ([]byte, error) {
	res, err := http.Post(url, "application/msgpack", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Server returned status code %v", res.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return bodyBytes
}

func postNoBody(url string) ([]byte, error) {
	var body []byte
	res, err := http.Post(url, "text/plain", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Server returned status code %v", res.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return bodyBytes
}

func get(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode !- http.StatusOK {
		return fmt.Errorf("Server returned status code %v", res.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return bodyBytes
}
