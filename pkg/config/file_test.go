package config

import (
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/tgoldstar/scm-poller/pkg/poller"
)

func TestFile_FetchNotExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	f := File{
		FileSystem: fs,
		Path:       "notexists.yaml",
	}
	opts, err := f.Fetch()

	assert.Nil(t, opts, "no opts should be returned")
	assert.Error(t, err, "fetch should fail")
}

func TestFile_FetchInvalidYAML(t *testing.T) {
	fs := afero.NewMemMapFs()
	fs.Mkdir("conf", 0755)
	afero.WriteFile(fs, "conf/poller.yaml", []byte("Invalid\nYAML\n"), 0644)
	f := File{
		FileSystem: fs,
		Path:       "conf/poller.yaml",
	}
	_, err := f.Fetch()

	assert.Error(t, err, "fetch should fail")
}

func TestFile_FetchInvalidFormat(t *testing.T) {
	conf := `
repositories:
- remote_url: "https://github.com/tgoldstar/scm-poller"
  reference: "master"
  interval: invalid
`
	fs := afero.NewMemMapFs()
	fs.Mkdir("conf", 0755)
	afero.WriteFile(fs, "conf/poller.yaml", []byte(conf), 0644)
	f := File{
		FileSystem: fs,
		Path:       "conf/poller.yaml",
	}
	_, err := f.Fetch()

	assert.Error(t, err, "fetch should fail")
}

func TestFile_FetchValidFormat(t *testing.T) {
	conf := `
repositories:
  - remote_url: "https://github.com/tgoldstar/scm-poller"
    reference: "master"
    interval: 3s
`
	fs := afero.NewMemMapFs()
	fs.Mkdir("conf", 0755)
	afero.WriteFile(fs, "conf/poller.yaml", []byte(conf), 0644)
	f := File{
		FileSystem: fs,
		Path:       "conf/poller.yaml",
	}
	opts, err := f.Fetch()

	assert.NotNil(t, opts, "initialized opts should be returned")
	assert.NoError(t, err, "fetch should not fail")
	assert.Len(t, opts.Repositories, 1, "one repository should be returned")
	assert.Equal(t, poller.RepositoryPollOpts{
		RemoteURL: "https://github.com/tgoldstar/scm-poller",
		Reference: "master",
		Interval:  time.Second * 3,
	}, opts.Repositories[0])
}
