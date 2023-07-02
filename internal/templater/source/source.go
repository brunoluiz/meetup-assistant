package source

import (
	"io/fs"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
)

// New parse the DSN and returns a FS
// DSN format: <srv>://<repo url>#<ref> (see github.com/hairyhenderson/go-fsimpl)
// Example git+file:///Users/bruno.silva/git/personal/templates
func New(dsn string) (fs.FS, error) {
	mux := fsimpl.NewMux()
	mux.Add(gitfs.FS)

	return mux.Lookup(dsn)
}
