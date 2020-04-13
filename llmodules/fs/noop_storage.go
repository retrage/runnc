package fs

import (
	"os"

	ll "github.com/retrage/runnc/llif"
)

type noopFsHandler struct{}

func NewNoopFsHandler() (ll.FsHandler, error) {
	return &noopFsHandler{}, nil
}

func (h *noopFsHandler) FsCreateFunc(i *ll.FsCreateInput) (*ll.LLState, error) {
	ret := &ll.LLState{}
	return ret, nil
}

func (h *noopFsHandler) FsRunFunc(i *ll.FsRunInput) (*ll.LLState, error) {
	return i.FsState, nil
}

func (h *noopFsHandler) FsDestroyFunc(i *ll.FsDestroyInput) (*ll.LLState, error) {
	if err := os.RemoveAll(i.ContainerRoot); err != nil {
		return nil, err
	}
	return i.FsState, nil
}
