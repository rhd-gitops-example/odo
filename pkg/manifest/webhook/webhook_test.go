package webhook

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/h2non/gock"

	"github.com/google/go-cmp/cmp"
)

var mockHeaders = map[string]string{
	"X-GitHub-Request-Id":   "DD0E:6011:12F21A8:1926790:5A2064E2",
	"X-RateLimit-Limit":     "60",
	"X-RateLimit-Remaining": "59",
	"X-RateLimit-Reset":     "1512076018",
}

func TestAddWebHook(t *testing.T) {

	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/foo/bar/hooks").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook.json")

	err := AddWebHook("https://github.com/foo/bar.git", "1013f24d30ba050e6f38be827c34b7da2682169b", "http://example.com/webhook", "mysecret")
	if err != nil {
		t.Fatal(err)
	}
}

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
