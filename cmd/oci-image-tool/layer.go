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
	"log"
	"os"

	"github.com/opencontainers/image-spec/image"
	"github.com/spf13/cobra"
)

type layerCmd struct {
	stdout *log.Logger
	stderr *log.Logger
}

func newLayerCmd(stdout, stderr *log.Logger) *cobra.Command {
	v := &layerCmd{
		stdout: stdout,
		stderr: stderr,
	}

	cmd := &cobra.Command{
		Use:   "create-layer [parent] [child]",
		Short: "Create an OCI layer",
		Long:  `Create an OCI layer based on the changeset between filesystems.`,
		Run:   v.Run,
	}

	return cmd
}

func (v *layerCmd) Run(cmd *cobra.Command, args []string) {
	if len(args) != 1 && len(args) != 2 {
		v.stderr.Print("One or two filesystems are required")
		if err := cmd.Usage(); err != nil {
			v.stderr.Println(err)
		}
		os.Exit(1)
	}
	var err error
	if len(args) == 1 {
		err = image.CreateLayer(args[0], "")
	} else {
		err = image.CreateLayer(args[0], args[1])
	}
	if err != nil {
		v.stderr.Printf("create layer failed: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
