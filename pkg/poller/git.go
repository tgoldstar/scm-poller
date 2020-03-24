package poller

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Git is a git client that implemnts the SourceControl interface
type Git struct {
}

var _ SourceControl = Git{}

// GetRef retrieves the git commit id of a reference in repository
func (g Git) GetRef(repository, reference string) (string, error) {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name:  "origin",
		URLs:  []string{repository},
		Fetch: []config.RefSpec{config.RefSpec(reference)},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "failed listing references on remote %s", repository)
	}

	for _, ref := range refs {
		if ref.Name().Short() == reference {
			return ref.Hash().String(), nil
		}
	}

	return "", errors.Errorf("reference %s not found in %s", reference, repository)
}
