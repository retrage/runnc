package main

import (
	"github.com/retrage/runnc/llcli"
	ll "github.com/retrage/runnc/llif"
	llfs "github.com/retrage/runnc/llmodules/fs"
	llnet "github.com/retrage/runnc/llmodules/network"
	llnabla "github.com/retrage/runnc/llruntimes/nabla"
)

func main() {
	fsH, err := llfs.NewExt4FsHandler()
	if err != nil {
		panic(err)
	}
	networkH, err := llnet.NewTapBrNetworkHandler()
	if err != nil {
		panic(err)
	}
	execH, err := llnabla.NewNablaExecHandler()
	if err != nil {
		panic(err)
	}

	nablaLLCHandler := ll.RunllcHandler{
		FsH:      fsH,
		NetworkH: networkH,
		ExecH:    execH,
	}

	// We run the OCI runtime called "runnc", with root dir "/run/runnc"
	// with the low level handlers chosen above.
	llcli.Runllc("runnc", "/run/runnc", nablaLLCHandler)
}
