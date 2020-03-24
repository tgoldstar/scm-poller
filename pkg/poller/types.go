package poller

import (
	"fmt"
	"time"
)

// RepositoryPollOpts defines configuration for a single reference polling
type RepositoryPollOpts struct {
	RemoteURL string `yaml:"remote_url"`
	Reference string
	Interval  time.Duration
}

// PollOpts defines configuration for polling,
// including what to poll and polling interval
type PollOpts struct {
	Repositories []RepositoryPollOpts
}

// Change is a difference in repository
type Change struct {
	Repository string
	Reference  string
	Old        string
	New        string
}

// String returns the string representation of Change
func (c Change) String() string {
	return fmt.Sprintf("Change in %s/%s (%s -> %s)", c.Repository, c.Reference, c.Old, c.New)
}

// PollError represents an error during polling
type PollError struct {
	Err        error
	Repository string
	Reference  string
}

// Error describes the polling error
func (p PollError) Error() string {
	return fmt.Sprintf("polling error at %s/%s: %v", p.Repository, p.Reference, p.Err)
}
