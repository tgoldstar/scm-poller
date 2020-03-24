package config

import "github.com/tgoldstar/scm-poller/pkg/poller"

// Fetcher is an abstraction around coniguration fetching
type Fetcher interface {
	Fetch() (*poller.PollOpts, error)
}
