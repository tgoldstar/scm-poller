package poller

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tgoldstar/scm-poller/pkg/store"
)

type FakeGit struct {
	mock.Mock
}

func (f *FakeGit) GetRef(repository, reference string) (string, error) {
	args := f.Called(repository, reference)
	return args.String(0), args.Error(1)
}

type FakeStore struct {
	mock.Mock
}

func (f *FakeStore) UpdateCommit(repo, ref, commit string) (string, error) {
	args := f.Called(repo, ref, commit)
	return args.String(0), args.Error(1)
}

var discardLogger = log.New(ioutil.Discard, "", 0)

func TestPoller_StartStop(t *testing.T) {
	p := New(discardLogger, store.NewMemoryStore(), &FakeGit{})
	c := make(chan Change, 1)
	e := make(chan PollError, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go cancel()
	p.Poll(ctx, &PollOpts{Repositories: []RepositoryPollOpts{}}, c, e)
	assert.Empty(t, c)
	assert.Empty(t, e)
}

func TestPoller_PollError(t *testing.T) {
	git := &FakeGit{}
	p := New(discardLogger, store.NewMemoryStore(), git)
	c := make(chan Change, 1)
	e := make(chan PollError, 1)
	ctx, cancel := context.WithCancel(context.Background())

	git.On("GetRef", "https://github.com/not/exists", "master").
		Return("", errors.New("repo not found"))

	go p.Poll(ctx, &PollOpts{Repositories: []RepositoryPollOpts{
		RepositoryPollOpts{
			RemoteURL: "https://github.com/not/exists",
			Reference: "master",
			Interval:  time.Millisecond * 10,
		},
	}}, c, e)

	select {
	case err := <-e:
		assert.Error(t, err)
	case <-time.NewTimer(time.Second * 1).C:
		t.Error("timed out, no error received")
	}

	cancel()
	assert.Empty(t, c, "No changes should be recorded")
}

func TestPoller_ReconcileRepoGetRefError(t *testing.T) {
	git := &FakeGit{}
	p := New(discardLogger, store.NewMemoryStore(), git).(*poller)
	c := make(chan Change, 1)

	git.On("GetRef", "https://github.com/not/exists", "master").
		Return("", errors.New("repo does not exist"))

	err := p.reconcile("https://github.com/not/exists", "master", c)

	assert.Empty(t, c, "no changes should be captured")
	assert.Error(t, err, "should fail by mock")
}

func TestPoller_ReconcileStoreError(t *testing.T) {
	git := &FakeGit{}
	st := &FakeStore{}
	p := New(discardLogger, st, git).(*poller)
	c := make(chan Change, 1)

	git.On("GetRef", "https://github.com/kubernetes/kubernetes", "master").
		Return("new", nil)
	st.On("UpdateCommit", "https://github.com/kubernetes/kubernetes", "master", "new").
		Return("", errors.New("update failed"))

	err := p.reconcile("https://github.com/kubernetes/kubernetes", "master", c)

	assert.Empty(t, c, "no changes should be captured")
	assert.Error(t, err, "should fail by mock")
}

func TestPoller_ReconcileSuccessNoChange(t *testing.T) {
	git := &FakeGit{}
	st := &FakeStore{}
	p := New(discardLogger, st, git).(*poller)
	c := make(chan Change, 1)

	git.On("GetRef", "https://github.com/kubernetes/kubernetes", "master").
		Return("old", nil)
	st.On("UpdateCommit", "https://github.com/kubernetes/kubernetes", "master", "old").
		Return("old", nil)

	err := p.reconcile("https://github.com/kubernetes/kubernetes", "master", c)

	assert.NoError(t, err)
	assert.Empty(t, c, "no changes should be captured")
}

func TestPoller_ReconcileSuccessChange(t *testing.T) {
	git := &FakeGit{}
	st := &FakeStore{}
	p := New(discardLogger, st, git).(*poller)
	c := make(chan Change, 1)

	git.On("GetRef", "https://github.com/kubernetes/kubernetes", "master").
		Return("new", nil)
	st.On("UpdateCommit", "https://github.com/kubernetes/kubernetes", "master", "new").
		Return("old", nil)

	err := p.reconcile("https://github.com/kubernetes/kubernetes", "master", c)

	assert.NoError(t, err)
	assert.NotEmpty(t, c)

	change := <-c

	assert.Equal(t, Change{
		Repository: "https://github.com/kubernetes/kubernetes",
		Reference:  "master",
		Old:        "old",
		New:        "new",
	}, change)
}
