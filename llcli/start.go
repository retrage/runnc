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

package llcli

import (
	"fmt"
	"github.com/retrage/runnc/libcontainer"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	ll "github.com/retrage/runnc/llif"
)

func newStartCmd(llcHandler ll.RunllcHandler, sf stringSubFunc) cli.Command {
	return cli.Command{
		Name:  "start",
		Usage: "executes the user defined process in a created container",
		ArgsUsage: sf(`<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`),
		Description: sf(`The start command executes the user defined process in a created container.`),
		Action: func(context *cli.Context) error {
			container, err := getContainer(context, llcHandler)
			if err != nil {
				return err
			}
			status, err := container.Status()
			if err != nil {
				return err
			}
			switch status {
			case libcontainer.Created:
				return container.Exec()
			case libcontainer.Stopped:
				return errors.New("cannot start a container that has stopped")
			case libcontainer.Running:
				return errors.New("cannot start an already running container")
			default:
				return fmt.Errorf("cannot start a container in the %s state\n", status)
			}
		},
	}
}
