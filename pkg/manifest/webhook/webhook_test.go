package webhook

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetDriverName(t *testing.T) {

	tests := []struct {
		url          string
		driver       string
		driverErrMsg string
		repo         string
		repoErrMsg   string
	}{
		{
			"http://github.org",
			"github",
			"",
			"",
			"failed to get Git repo: ",
		},
		{
			"http://github.com/",
			"github",
			"",
			"",
			"failed to get Git repo: /",
		},
		{
			"http://github.com/foo/bar",
			"github",
			"",
			"foo/bar",
			"",
		},
		{
			"https://githuB.com/foo/bar.git",
			"github",
			"",
			"foo/bar",
			"",
		},
		{
			"http://gitlab.com/foo/bar.git2",
			"gitlab",
			"",
			"",
			"failed to get Git repo: /foo/bar.git2",
		},
		{
			"http://gitlab/foo/bar/",
			"",
			"unknown Git server: gitlab",
			"foo/bar",
			"",
		},
		{
			"https://gitlab.a.b/foo/bar/bar",
			"",
			"unknown Git server: gitlab.a.b",
			"",
			"failed to get Git repo: /foo/bar/bar",
		},
		{
			"https://gitlab.org2/f.b/bar.git",
			"",
			"unknown Git server: gitlab.org2",
			"",
			"failed to get Git repo: /f.b/bar.git",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			u, err := url.Parse(test.url)
			if err != nil {
				t.Error(err)
			} else {
				gotDriver, err := getDriverName(u)
				driverErrMsg := ""
				if err != nil {
					driverErrMsg = err.Error()
				}
				if diff := cmp.Diff(test.driverErrMsg, driverErrMsg); diff != "" {
					t.Errorf("driver errMsg mismatch got\n%s", diff)
				}
				if diff := cmp.Diff(test.driver, gotDriver); diff != "" {
					t.Errorf("driver mismatch got\n%s", diff)
				}

				repo, err := getRepo(u)
				repoErrMsg := ""
				if err != nil {
					repoErrMsg = err.Error()
				}
				if diff := cmp.Diff(test.repoErrMsg, repoErrMsg); diff != "" {
					t.Errorf("driver errMsg mismatch got\n%s", diff)
				}
				if diff := cmp.Diff(test.repo, repo); diff != "" {
					t.Errorf("driver mismatch got\n%s", diff)
				}

			}
		})
	}
}
