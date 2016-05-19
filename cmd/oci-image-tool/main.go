package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "oci-image-tool",
		Short: "A tool for working with OCI images",
	}

	stdout := log.New(os.Stdout, "", 0)
	stderr := log.New(os.Stderr, "", 0)

	cmd.AddCommand(newValidateCmd(stdout, stderr))
	if err := cmd.Execute(); err != nil {
		stderr.Println(err)
		os.Exit(1)
	}
}
