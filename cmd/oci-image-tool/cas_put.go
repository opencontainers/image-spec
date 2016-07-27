// Copyright 2016 The Linux Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/opencontainers/image-spec/image/cas/layout"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type casPutCmd struct {
	stdout io.Writer
	stderr *log.Logger
	path   string
}

func newCASPutCmd(stdout io.Writer, stderr *log.Logger) *cobra.Command {
	state := &casPutCmd{
		stdout: stdout,
		stderr: stderr,
	}

	return &cobra.Command{
		Use:   "put PATH",
		Short: "Write a blob to the store",
		Long:  "Read a blob from stdin, write it to the store, and print the digest to stdout.",
		Run:   state.Run,
	}
}

func (state *casPutCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		if err := cmd.Usage(); err != nil {
			state.stderr.Println(err)
		}
		os.Exit(1)
	}

	state.path = args[0]

	err := state.run()
	if err != nil {
		state.stderr.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func (state *casPutCmd) run() (err error) {
	ctx := context.Background()

	engine, err := layout.NewEngine(ctx, state.path)
	if err != nil {
		return err
	}
	defer engine.Close()

	digest, err := engine.Put(ctx, os.Stdin)
	if err != nil {
		return err
	}

	n, err := fmt.Fprintln(state.stdout, digest)
	if err != nil {
		return err
	}
	if n < len(digest) {
		return fmt.Errorf("wrote %d of %d bytes", n, len(digest))
	}

	return nil
}
