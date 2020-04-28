package webhook

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetDriverName(t *testing.T) {

	tests := []struct {
		url    string
		want   string
		errMsg string
	}{
		{
			"http://github.org",
			"github",
			"",
		},
		{
			"http://github.com",
			"github",
			"",
		},
		{
			"http://github.com/foo/bar",
			"github",
			"",
		},
		{
			"https://githuB.com/foo/bar",
			"github",
			"",
		},
		{
			"http://gitlab.com/foo/bar",
			"gitlab",
			"",
		},
		{
			"http://gitlab/foo/bar",
			"",
			"unknown Git server: gitlab",
		},
		{
			"https://gitlab.a.b/foo/bar",
			"",
			"unknown Git server: gitlab.a.b",
		},
	}

	i := 0
	for _, test := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			u, err := url.Parse(test.url)
			if err != nil {
				t.Error(err)
			} else {
				got, err := getDriverName(u)
				errMsg := ""
				if err != nil {
					errMsg = err.Error()
				}
				if diff := cmp.Diff(test.errMsg, errMsg); diff != "" {
					t.Errorf("errMsg mismatch got\n%s", diff)
				}
				if diff := cmp.Diff(test.want, got); diff != "" {
					t.Errorf("driver mismatch got\n%s", diff)
				}
			}
		})
		i++
	}
}
