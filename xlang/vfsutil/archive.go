package vfsutil

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/sourcegraph/sourcegraph/pkg/gitserver"
	"github.com/sourcegraph/sourcegraph/pkg/vcs/git"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sourcegraph/ctxvfs"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

// ArchiveFileSystem returns a virtual file system backed by a .zip
// archive of a Git tree (in the common case, the root tree of a Git
// repository at a specific commit). The treeish is a Git object ID
// that refers to a tree; it can be a commit ID, a tree ID, and so
// on. For consistency, callers should generally use full SHAs, not
// rev specs like branch names, etc.
//
// ArchiveFileSystem fetches the full .zip archive initially and then
// can satisfy FS operations nearly instantly in memory.
func ArchiveFileSystem(repo gitserver.Repo, treeish string) *ArchiveFS {
	fetch := func(ctx context.Context) (*archiveReader, error) {
		rc, err := git.Archive(ctx, repo, git.ArchiveOptions{
			Treeish: treeish,
			Format:  "zip",
		})
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		data, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil, err
		}
		gitserverBytes.Add(float64(len(data)))

		zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return nil, err
		}
		return &archiveReader{
			Reader: zr,
		}, nil
	}
	return &ArchiveFS{fetch: fetch}
}

// archiveReader is like zip.ReadCloser, but it allows us to use a custom
// closer.
type archiveReader struct {
	*zip.Reader
	io.Closer
	Evicter

	// Prefix is the path prefix to strip. For example a GitHub archive
	// has a top-level dir "{repobasename}-{sha}/".
	Prefix string
}

// ArchiveFS is a ctxvfs.FileSystem backed by an Archiver.
type ArchiveFS struct {
	fetch func(context.Context) (*archiveReader, error)

	// EvictOnClose when true will evict the underlying archive from the
	// archive cache when closed.
	EvictOnClose bool

	once sync.Once
	err  error // the error encountered during the fetch call (if any)
	ar   *archiveReader
	fs   vfs.FileSystem // the zipfs virtual file system

	// We have a mutex for closed to prevent Close and fetch racing.
	closedMu sync.Mutex
	closed   bool
}

// fetchOrWait initiates the fetch if it has not yet
// started. Otherwise it waits for it to finish.
func (fs *ArchiveFS) fetchOrWait(ctx context.Context) error {
	fs.once.Do(func() {
		// If we have already closed, do not open new resources. If we
		// haven't closed, prevent closing while fetching by holding
		// the lock.
		fs.closedMu.Lock()
		defer fs.closedMu.Unlock()
		if fs.closed {
			fs.err = errors.New("closed")
			return
		}

		fs.ar, fs.err = fs.fetch(ctx)
		if fs.err == nil {
			fs.fs = zipfs.New(&zip.ReadCloser{Reader: *fs.ar.Reader}, "")
			if fs.ar.Prefix != "" {
				ns := vfs.NameSpace{}
				ns.Bind("/", fs.fs, "/"+fs.ar.Prefix, vfs.BindReplace)
				fs.fs = ns
			}
		}
	})
	return fs.err
}

func (fs *ArchiveFS) Open(ctx context.Context, name string) (ctxvfs.ReadSeekCloser, error) {
	if err := fs.fetchOrWait(ctx); err != nil {
		return nil, err
	}
	return fs.fs.Open(name)
}

func (fs *ArchiveFS) Lstat(ctx context.Context, path string) (os.FileInfo, error) {
	if err := fs.fetchOrWait(ctx); err != nil {
		return nil, err
	}
	return fs.fs.Lstat(path)
}

func (fs *ArchiveFS) Stat(ctx context.Context, path string) (os.FileInfo, error) {
	if err := fs.fetchOrWait(ctx); err != nil {
		return nil, err
	}
	return fs.fs.Stat(path)
}

func (fs *ArchiveFS) ReadDir(ctx context.Context, path string) ([]os.FileInfo, error) {
	if err := fs.fetchOrWait(ctx); err != nil {
		return nil, err
	}
	return fs.fs.ReadDir(path)
}

func (fs *ArchiveFS) ListAllFiles(ctx context.Context) ([]string, error) {
	if err := fs.fetchOrWait(ctx); err != nil {
		return nil, err
	}

	filenames := make([]string, 0, len(fs.ar.File))
	for _, f := range fs.ar.File {
		if f.Mode().IsRegular() {
			filenames = append(filenames, strings.TrimPrefix(f.Name, fs.ar.Prefix))
		}
	}
	return filenames, nil
}

func (fs *ArchiveFS) Close() error {
	fs.closedMu.Lock()
	defer fs.closedMu.Unlock()
	if fs.closed {
		return errors.New("already closed")
	}

	fs.closed = true
	if fs.ar != nil && fs.ar.Closer != nil {
		err := fs.ar.Close()
		if err != nil {
			return err
		}
		if fs.EvictOnClose && fs.ar.Evicter != nil {
			fs.ar.Evict()
		}
	}
	return nil
}

func (fs *ArchiveFS) String() string { return "ArchiveFS(" + fs.fs.String() + ")" }

var gitserverBytes = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "src",
	Subsystem: "vfs",
	Name:      "gitserver_bytes_total",
	Help:      "Total number of bytes read into memory by ArchiveFileSystem.",
})

func init() {
	prometheus.MustRegister(gitserverBytes)
}
