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

	"github.com/opencontainers/image-spec/image/refs/layout"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type refsListCmd struct {
	stdout io.Writer
	stderr *log.Logger
	path   string
}

func newRefsListCmd(stdout io.Writer, stderr *log.Logger) *cobra.Command {
	state := &refsListCmd{
		stdout: stdout,
		stderr: stderr,
	}

	return &cobra.Command{
		Use:   "list PATH",
		Short: "Return available names from the store.",
		Run:   state.Run,
	}
}

func (state *refsListCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		state.stderr.Print("PATH must be provided")
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

func (state *refsListCmd) run() (err error) {
	ctx := context.Background()

	engine, err := layout.NewEngine(state.path)
	if err != nil {
		return err
	}
	defer engine.Close()

	return engine.List(ctx, "", -1, 0, state.printName)
}

func (state *refsListCmd) printName(ctx context.Context, name string) (err error) {
	n, err := io.WriteString(state.stdout, fmt.Sprintf("%s\n", name))
	if err != nil {
		return err
	}
	if n < len(name)+1 {
		err = fmt.Errorf("wrote %d of %d characters", n, len(name)+1)
		return err
	}
	return nil
}
