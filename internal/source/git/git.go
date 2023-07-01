package git

import (
	"io/fs"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
)

func New(address string) (fs.FS, error) {
	mux := fsimpl.NewMux()
	mux.Add(gitfs.FS)

	return mux.Lookup(address)
}
