package source

import (
	"io/fs"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
)

func New(dsn string) (fs.FS, error) {
	mux := fsimpl.NewMux()
	mux.Add(gitfs.FS)

	return mux.Lookup(dsn)
}
