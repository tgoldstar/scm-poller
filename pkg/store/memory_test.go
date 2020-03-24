package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStore_CreateRepository(t *testing.T) {
	m := &Memory{
		data: make(map[string]repository),
	}
	old, err := m.UpdateCommit("github.com/tgoldstar/scm-poller", "master", "somecommit")

	assert.NoError(t, err, "memory store should not fail")
	assert.Equal(t, "", old, "old should be empty")
	assert.Equal(t, repository{
		"master": "somecommit",
	}, m.data["github.com/tgoldstar/scm-poller"])
}

func TestMemoryStore_CreateReference(t *testing.T) {
	m := &Memory{
		data: make(map[string]repository),
	}
	m.UpdateCommit("github.com/tgoldstar/scm-poller", "master", "somecommit")
	old, err := m.UpdateCommit("github.com/tgoldstar/scm-poller", "stable", "othercommit")

	assert.Equal(t, "", old, "old should be empty")
	assert.NoError(t, err, "memory store should not fail")
	assert.Equal(t, repository{
		"master": "somecommit",
		"stable": "othercommit",
	}, m.data["github.com/tgoldstar/scm-poller"])
}

func TestMemoryStore_UpdateCommit(t *testing.T) {
	m := &Memory{
		data: make(map[string]repository),
	}
	m.UpdateCommit("github.com/tgoldstar/scm-poller", "master", "somecommit")
	old, err := m.UpdateCommit("github.com/tgoldstar/scm-poller", "master", "othercommit")

	assert.Equal(t, "somecommit", old)
	assert.NoError(t, err, "memory store should not fail")
	assert.Equal(t, repository{
		"master": "othercommit",
	}, m.data["github.com/tgoldstar/scm-poller"])
}
