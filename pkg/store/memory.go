package store

// repository maps between branch name and last commit
type repository map[string]string

// Memory is an in-memory store
type Memory struct {
	data map[string]repository
}

// NewMemoryStore initializes an empty in-memory store
func NewMemoryStore() Store {
	return &Memory{
		data: make(map[string]repository),
	}
}

// UpdateCommit updates the last known commit id of a branch
func (m *Memory) UpdateCommit(remote, branch, commit string) (string, error) {
	branches, ok := m.data[remote]
	var previous string

	if ok {
		previous = branches[branch]
	} else {
		m.data[remote] = make(repository)
		branches = m.data[remote]
	}

	branches[branch] = commit

	return previous, nil
}
