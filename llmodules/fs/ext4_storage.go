package fs

import (
	"os"
	"path/filepath"

	"github.com/retrage/runnc/libcontainer/configs"
	ll "github.com/retrage/runnc/llif"
	"github.com/retrage/runnc/nabla-lib/storage"
	"github.com/retrage/runnc/utils"
	"github.com/pkg/errors"
)

type Ext4FsHandler struct{}

func NewExt4FsHandler() (ll.FsHandler, error) {
	return &Ext4FsHandler{}, nil
}

func (h *Ext4FsHandler) FsCreateFunc(i *ll.FsCreateInput) (*ll.LLState, error) {
	fsPath, err := createRootfsExt4(i.Config, i.ContainerRoot)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create rootfs disk image")
	}

	ret := &ll.LLState{}
	ret.Options = map[string]string{
		"FsPath": fsPath,
	}

	return ret, nil
}

func (h *Ext4FsHandler) FsRunFunc(i *ll.FsRunInput) (*ll.LLState, error) {
	return i.FsState, nil
}

func (h *Ext4FsHandler) FsDestroyFunc(i *ll.FsDestroyInput) (*ll.LLState, error) {
	if err := os.RemoveAll(i.ContainerRoot); err != nil {
		return nil, err
	}
	return i.FsState, nil
}

func createRootfsExt4(config *configs.Config, containerRoot string) (string, error) {
	rootfsPath := config.Rootfs
	targetExt4Path := filepath.Join(containerRoot, "rootfs.img")
	if err := os.MkdirAll(filepath.Join(rootfsPath, "/etc"), 0755); err != nil {
		return "", errors.Wrap(err, "Unable to create "+filepath.Join(rootfsPath, "/etc"))
	}
	for _, mount := range config.Mounts {
		if (mount.Destination == "/etc/resolv.conf") ||
			(mount.Destination == "/etc/hosts") ||
			(mount.Destination == "/etc/hostname") {
			dest := filepath.Join(rootfsPath, mount.Destination)
			source := mount.Source
			if err := utils.Copy(dest, source); err != nil {
				return "", errors.Wrap(err, "Unable to copy "+source+" to "+dest)
			}
		}
	}
	_, err := storage.CreateExt4(rootfsPath, &targetExt4Path)
	if err != nil {
		return "", errors.Wrap(err, "Error creating disk image from rootfs")
	}
	return targetExt4Path, nil
}
