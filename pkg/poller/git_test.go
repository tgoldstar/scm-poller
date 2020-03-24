package poller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit_GetRefRepoNotExists(t *testing.T) {
	g := &Git{}
	_, err := g.GetRef("https://github.com/tgoldstar/notexists", "master")
	assert.Error(t, err, "repository does not exist, should fail")
}

func TestGit_GetRefRefNotExists(t *testing.T) {
	g := &Git{}
	_, err := g.GetRef("https://github.com/tgoldstar/scm-poller", "notexists")
	assert.Error(t, err, "reference does not exist, should fail")
}

func TestGit_GetRefSuccess(t *testing.T) {
	g := &Git{}
	c, err := g.GetRef("https://github.com/aquasecurity/kube-hunter", "v0.2.0")

	assert.NoError(t, err, "reference exists, should not fail")
	assert.Equal(t, "1d7bdd613109cde1134133b2560b92817432b530", c)
}
