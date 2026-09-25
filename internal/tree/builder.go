package tree

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"syscall"
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

	r := new(runner{
		ctx:        ctx,
		cancel:     cancel,
		sem:        make(chan struct{}, concurrency),
		permission: opts.Permission,
		gitkeep:    opts.Gitkeep,
	})

	r.processNodes(nodes, targetDir)
	r.wg.Wait()

	if r.firstErr != nil {
		return r.firstErr
	}

	return nil
}

// createEmptyFile creates an empty file using low-level syscalls to eliminate
// runtime poller registration, finalizer, and heap allocation overheads of os.Create.
func createEmptyFile(path string) error {
	fd, err := syscall.Open(path, syscall.O_CREAT|syscall.O_WRONLY|syscall.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	return syscall.Close(fd)
}

func joinPath(dir, name string) string {
	if dir == "." {
		return name
	}
	if len(dir) > 0 && (dir[len(dir)-1] == '/' || dir[len(dir)-1] == filepath.Separator) {
		return dir + name
	}
	return dir + string(filepath.Separator) + name
}

func buildSynchronous(nodes []Node, parentDir string, permission os.FileMode, gitkeep bool) error {
	for i := range nodes {
		if nodes[i].Type != NodeTypeDirectory {
			continue
		}
		path := joinPath(parentDir, nodes[i].Name)
		if err := os.Mkdir(path, permission); err != nil {
			return err
		}

		if gitkeep {
			if err := createEmptyFile(joinPath(path, ".gitkeep")); err != nil {
				return err
			}
		}

		if len(nodes[i].Contents) > 0 {
			if err := buildSynchronous(nodes[i].Contents, path, permission, gitkeep); err != nil {
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
	// Count directories without allocating a new slice
	var dirIndices []int
	for i := range nodes {
		if nodes[i].Type == NodeTypeDirectory {
			dirIndices = append(dirIndices, i)
		}
	}

	if len(dirIndices) == 0 {
		return
	}

	// Run all but the last directory in worker goroutines if semaphore has space,
	// and run the last directory directly in the current goroutine.
	for i := 0; i < len(dirIndices)-1; i++ {
		if r.ctx.Err() != nil {
			return
		}

		nodeIdx := dirIndices[i]
		select {
		case r.sem <- struct{}{}:
			r.wg.Add(1)
			go func(idx int) {
				defer func() {
					<-r.sem
					r.wg.Done()
				}()
				r.processSingleDir(&nodes[idx], parentDir)
			}(nodeIdx)
		case <-r.ctx.Done():
			return
		default:
			// Run synchronously if semaphore is occupied
			r.processSingleDir(&nodes[nodeIdx], parentDir)
		}
	}

	// Execute the final directory in the current goroutine to save goroutine spawn costs
	lastIdx := dirIndices[len(dirIndices)-1]
	r.processSingleDir(&nodes[lastIdx], parentDir)
}

func (r *runner) processSingleDir(n *Node, parentDir string) {
	if r.ctx.Err() != nil {
		return
	}

	path := joinPath(parentDir, n.Name)
	if err := os.Mkdir(path, r.permission); err != nil {
		r.setError(err)
		return
	}

	if r.gitkeep {
		if err := createEmptyFile(joinPath(path, ".gitkeep")); err != nil {
			r.setError(err)
			return
		}
	}

	if len(n.Contents) > 0 {
		r.processNodes(n.Contents, path)
	}
}
