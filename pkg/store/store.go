package store

// Store is an interface for fetching data from external source
type Store interface {
	// UpdateCommit updates the last known commit id of a branch
	// and returns the previous one
	UpdateCommit(remote, branch, commit string) (string, error)
}
