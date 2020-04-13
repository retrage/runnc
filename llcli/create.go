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
	"os"

	ll "github.com/retrage/runnc/llif"
	"github.com/urfave/cli"
)

func newCreateCmd(llcHandler ll.RunllcHandler, sf stringSubFunc) cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "create a container",
		ArgsUsage: sf(`<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`),
		Description: sf(`The create command creates an instance of a container for a bundle. The bundle
is a directory with a specification file named "` + specConfig + `" and a root
filesystem.

The specification file includes an args parameter. The args parameter is used
to specify command(s) that get run when the container is started.
`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "bundle, b",
				Value: "",
				Usage: `path to the root of the bundle directory, defaults to the current directory`,
			},
			cli.StringFlag{
				Name:  "console-socket",
				Value: "",
				Usage: "path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal",
			},
			cli.StringFlag{
				Name:  "pid-file",
				Value: "",
				Usage: "specify the file to write the process id to",
			},
			cli.BoolFlag{
				Name:  "no-pivot",
				Usage: "do not use pivot root to jail process inside rootfs.  This should be used whenever the rootfs is on top of a ramdisk",
			},
			cli.BoolFlag{
				Name:  "no-new-keyring",
				Usage: "do not create a new session keyring for the container.  This will cause the container to inherit the calling processes session key",
			},
			cli.IntFlag{
				Name:  "preserve-fds",
				Usage: "Pass N additional file descriptors to the container (stdio + $LISTEN_FDS + N in total)",
			},
		},
		Action: func(context *cli.Context) error {
			// TODO: Implement
			spec, err := setupSpec(context)
			if err != nil {
				fatal(err)
			}

			status, err := startContainer(context, llcHandler, spec, true)
			if err != nil {
				fatal(err)
			}

			// exit with the container's exit status so any external supervisor is
			// notified of the exit with the correct exit status.
			os.Exit(status)
			return nil
		},
	}

}
