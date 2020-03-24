package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/tgoldstar/scm-poller/pkg/poller"
	"gopkg.in/yaml.v2"
)

var _ Fetcher = File{}

// File fetches configuration from YAML file
type File struct {
	FileSystem afero.Fs
	Path       string
}

// Fetch loads configuration from f.Path
func (f File) Fetch() (*poller.PollOpts, error) {
	content, err := afero.ReadFile(f.FileSystem, f.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading %s", f.Path)
	}

	var opts poller.PollOpts
	err = yaml.UnmarshalStrict(content, &opts)

	return &opts, errors.Wrap(err, "failed parsing configuration file")
}
