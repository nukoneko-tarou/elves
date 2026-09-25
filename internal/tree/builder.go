package tree

import (
	"context"
	"os"
	"path/filepath"
	"sync"
)

// DefaultConcurrency is the empirically optimal concurrency level for local filesystems,
// balancing I/O parallelism against OS kernel VFS/directory lock contention.
const DefaultConcurrency = 2

// BuildOptions defines configuration for directory tree construction.
type BuildOptions struct {
	BaseDir     string
	SubDir      string
	Permission  os.FileMode
	Gitkeep     bool
	Concurrency int
}

// Build creates the directory tree on the filesystem according to options.
func Build(nodes []Node, opts BuildOptions) error {
	baseDir := opts.BaseDir
	if baseDir == "" {
		baseDir = "."
	}

	targetDir := baseDir
	if opts.SubDir != "" {
		targetDir = filepath.Join(baseDir, opts.SubDir)
		if err := os.Mkdir(targetDir, opts.Permission); err != nil {
			return err
		}
	}

	concurrency := opts.Concurrency
	if concurrency <= 0 {
		concurrency = DefaultConcurrency
	}

	// For concurrency == 1, execute synchronously without goroutine or channel overhead
	if concurrency == 1 {
		return buildSynchronous(nodes, targetDir, opts.Permission, opts.Gitkeep)
	}

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	r := &runner{
		ctx:        ctx,
		cancel:     cancel,
		sem:        make(chan struct{}, concurrency),
		permission: opts.Permission,
		gitkeep:    opts.Gitkeep,
	}

	r.processNodes(nodes, targetDir)
	r.wg.Wait()

	if r.firstErr != nil {
		return r.firstErr
	}

	return nil
}

func buildSynchronous(nodes []Node, parentDir string, permission os.FileMode, gitkeep bool) error {
	dirs := filterDirectories(nodes)
	for _, n := range dirs {
		path := filepath.Join(parentDir, n.Name)
		if err := os.Mkdir(path, permission); err != nil {
			return err
		}

		if gitkeep {
			f, err := os.Create(filepath.Join(path, ".gitkeep"))
			if err != nil {
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		}

		if len(n.Contents) > 0 {
			if err := buildSynchronous(n.Contents, path, permission, gitkeep); err != nil {
				return err
			}
		}
	}
	return nil
}

type runner struct {
	ctx        context.Context
	cancel     context.CancelCauseFunc
	sem        chan struct{}
	wg         sync.WaitGroup
	permission os.FileMode
	gitkeep    bool
	errOnce    sync.Once
	firstErr   error
}

func (r *runner) setError(err error) {
	if err != nil {
		r.errOnce.Do(func() {
			r.firstErr = err
			r.cancel(err)
		})
	}
}

func (r *runner) processNodes(nodes []Node, parentDir string) {
	dirs := filterDirectories(nodes)
	for _, n := range dirs {
		if r.ctx.Err() != nil {
			return
		}

		childNode := n
		select {
		case r.sem <- struct{}{}:
			r.wg.Add(1)
			go func() {
				defer func() {
					<-r.sem
					r.wg.Done()
				}()
				r.processSingleDir(childNode, parentDir)
			}()
		case <-r.ctx.Done():
			return
		default:
			// When semaphore buffer is busy, execute synchronously in the current goroutine
			r.processSingleDir(childNode, parentDir)
		}
	}
}

func (r *runner) processSingleDir(n Node, parentDir string) {
	if r.ctx.Err() != nil {
		return
	}

	path := filepath.Join(parentDir, n.Name)
	if err := os.Mkdir(path, r.permission); err != nil {
		r.setError(err)
		return
	}

	if r.gitkeep {
		f, err := os.Create(filepath.Join(path, ".gitkeep"))
		if err != nil {
			r.setError(err)
			return
		}
		if err := f.Close(); err != nil {
			r.setError(err)
			return
		}
	}

	if len(n.Contents) > 0 {
		r.processNodes(n.Contents, path)
	}
}

func filterDirectories(nodes []Node) []Node {
	var dirs []Node
	for _, n := range nodes {
		if n.Type == NodeTypeDirectory {
			dirs = append(dirs, n)
		}
	}
	return dirs
}
