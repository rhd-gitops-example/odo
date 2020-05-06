package service

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/spf13/cobra"
)

type keyValuePair struct {
	key   string
	value string
}

func TestAddCommandWithMissingParams(t *testing.T) {

	// manifestFile := "~/pipelines.yaml"
	cmdTests := []struct {
		desc    string
		flags   []keyValuePair
		wantErr string
	}{
		{"Missing app-name flag",
			[]keyValuePair{
				flag("service-name", "sample"), flag("git-repo-url", "example/repo"), flag("webhook-secret", "abc123"), flag("env-name", "test")},
			`Required flag(s) "app-name" have/has not been set`},
		{"Missing service-name flag",
			[]keyValuePair{flag("app-name", "app"),
				flag("git-repo-url", "example/repo"), flag("webhook-secret", "abc123"), flag("env-name", "test")},
			`Required flag(s) "service-name" have/has not been set`},
		{"Missing env-name flag",
			[]keyValuePair{flag("app-name", "app"),
				flag("service-name", "sample"), flag("git-repo-url", "sample/repo"), flag("webhook-secret", "abc123")},
			`Required flag(s) "env-name" have/has not been set`},
	}
	for _, tt := range cmdTests {
		t.Run(tt.desc, func(t *testing.T) {
			_, _, err := executeCommand(newCmdAdd("add", "odo pipelines service"), tt.flags...)
			if err.Error() != tt.wantErr {
				t.Errorf("got %s, want %s", err, tt.wantErr)
			}
		})
	}
}

func executeCommand(cmd *cobra.Command, flags ...keyValuePair) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	for _, flag := range flags {
		cmd.Flags().Set(flag.key, flag.value)
	}
	c, err = cmd.ExecuteC()
	return c, buf.String(), err
}

func matchError(t *testing.T, s string, e error) bool {
	t.Helper()
	if s == "" && e == nil {
		return true
	}
	if s != "" && e == nil {
		return false
	}
	match, err := regexp.MatchString(s, e.Error())
	if err != nil {
		t.Fatal(err)
	}
	return match
}

func flag(k, v string) keyValuePair {
	return keyValuePair{
		key:   k,
		value: v,
	}
}
