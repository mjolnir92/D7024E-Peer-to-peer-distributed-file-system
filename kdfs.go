package main

import (
  "os"
	"github.com/mjolnir92/kdfs/cmd"
)

func main() {
  if err := cmd.RootCmd.Execute(); err != nil {
    os.Exit(1)
  }
}
