// Copyright 2014 Docker, Inc.
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

// +build linux

package llcli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/retrage/runnc/libcontainer"
	ll "github.com/retrage/runnc/llif"
	"github.com/urfave/cli"
)

func init() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()
	}
}

func newInitCmd(llcHandler ll.RunllcHandler, sf stringSubFunc) cli.Command {
	return cli.Command{
		Name:  "init",
		Usage: sf(`initialize the namespaces and launch the process (do not call it outside of {{name}})`),
		Action: func(context *cli.Context) error {
			factory, _ := libcontainer.New("", llcHandler)
			if err := factory.StartInitialization(); err != nil {
				// as the error is sent back to the parent there is no need to log
				// or write it to stderr because the parent process will handle this
				fmt.Fprintf(os.Stderr, "ERR: %v", err)
				fmt.Fprintf(os.Stdout, "ERR: %v", err)
				os.Exit(1)
			}
			panic("libcontainer: container init failed to exec")
		},
	}
}
