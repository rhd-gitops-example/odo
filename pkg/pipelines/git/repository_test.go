package git

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

var mockHeaders = map[string]string{
	"X-GitHub-Request-Id":   "DD0E:6011:12F21A8:1926790:5A2064E2",
	"X-RateLimit-Limit":     "60",
	"X-RateLimit-Remaining": "59",
	"X-RateLimit-Reset":     "1512076018",
}

func TestListWebHooks(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/foo/bar/hooks").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hooks.json")

	repo, err := newRepository("https://github.com/foo/bar.git", "token")
	if err != nil {
		t.Fatal(err)
	}

	ids, err := repo.listWebhooks("http://example.com/webhook")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(ids, []string{"1"}); diff != "" {
		t.Errorf("driver errMsg mismatch got\n%s", diff)
	}
}

func TestDeleteWebHooks(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Delete("/repos/foo/bar/hooks/1").
		Reply(204).
		Type("application/json").
		SetHeaders(mockHeaders)

	repo, err := newRepository("https://github.com/foo/bar.git", "token")
	if err != nil {
		t.Fatal(err)
	}

	err = repo.deleteWebhooks("http://example.com", []string{"1"})
	if err != nil {
		t.Fatal(err)
	}

}

func TestCreateWebHook(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/foo/bar/hooks").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/hook.json")

	repo, err := newRepository("https://github.com/foo/bar.git", "token")
	if err != nil {
		t.Fatal(err)
	}

	err = repo.createWehoook()"http://example.com/webhook", "mysecret")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDriverName(t *testing.T) {

	tests := []struct {
		url          string
		driver       string
		driverErrMsg string
		repoName     string
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

				repoName, err := getRepoName(u)
				repoErrMsg := ""
				if err != nil {
					repoErrMsg = err.Error()
				}
				if diff := cmp.Diff(test.repoErrMsg, repoErrMsg); diff != "" {
					t.Errorf("driver errMsg mismatch got\n%s", diff)
				}
				if diff := cmp.Diff(test.repoName, repoName); diff != "" {
					t.Errorf("driver mismatch got\n%s", diff)
				}

			}
		})
	}
}
