package repository

import (
	"github.com/openshift/odo/pkg/odo/cli/pipelines/scm"
	"github.com/openshift/odo/pkg/odo/cli/pipelines/scm/github"
	"github.com/openshift/odo/pkg/odo/cli/pipelines/scm/utility"
)

func New(rawURL string) (scm.Repository, error) {
	repoType, err := utility.GetDriverName(rawURL)
	if err != nil {
		return nil, err
	}
	switch repoType {
	case "github":
		return github.NewRepository(rawURL)
	}
	return nil, nil
}
