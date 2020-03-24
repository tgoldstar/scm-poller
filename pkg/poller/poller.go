package poller

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/tgoldstar/scm-poller/pkg/store"
)

// Poller is an abstract interface for polling
type Poller interface {
	// Poll checks for changes in a remote
	Poll(ctx context.Context, opts *PollOpts, changes chan<- Change, errs chan<- PollError)
}

// SourceControl is an abstract type of SCM
type SourceControl interface {
	GetRef(repository, reference string) (string, error)
}

type poller struct {
	logger *log.Logger
	store  store.Store
	scm    SourceControl
}

// New creates a new poller object
func New(logger *log.Logger, store store.Store, scm SourceControl) Poller {
	return &poller{
		logger: logger,
		store:  store,
		scm:    scm,
	}
}

// Poll checks for changes in a repository
func (p *poller) Poll(ctx context.Context, opts *PollOpts, changes chan<- Change, errors chan<- PollError) {
	var wg sync.WaitGroup

	p.logger.Printf("Polling %d remotes for changes...", len(opts.Repositories))
	for _, repo := range opts.Repositories {
		wg.Add(1)
		go p.pollSingleReference(ctx, repo, changes, errors, &wg)
	}

	<-ctx.Done()
	wg.Wait()
}

func (p *poller) pollSingleReference(ctx context.Context, opts RepositoryPollOpts, changes chan<- Change, errs chan<- PollError, wg *sync.WaitGroup) {
	p.logger.Printf("Starting to poll %s (%s) every %v", opts.RemoteURL, opts.Reference, opts.Interval)
	pollTicker := time.NewTicker(opts.Interval)

	for {
		select {
		case <-pollTicker.C:
			err := p.reconcile(opts.RemoteURL, opts.Reference, changes)
			if err != nil {
				errs <- PollError{
					Err:        err,
					Repository: opts.RemoteURL,
					Reference:  opts.Reference,
				}
			}
		case <-ctx.Done():
			wg.Done()
			return
		}
	}
}

func (p *poller) reconcile(remote, reference string, changes chan<- Change) error {
	current, err := p.scm.GetRef(remote, reference)
	if err != nil {
		return err
	}

	previous, err := p.store.UpdateCommit(remote, reference, current)
	if err != nil {
		return errors.Wrap(err, "failed to load from store")
	}

	if current != previous {
		p.logger.Printf("Updated %s/%s to %s", remote, reference, current)
		changes <- Change{
			Repository: remote,
			Reference:  reference,
			New:        current,
			Old:        previous,
		}
	}

	return nil
}
