package pipelines

import (
	"bytes"
	"regexp"
	"testing"

	util "github.com/openshift/odo/pkg/odo/util"
	"github.com/spf13/cobra"
)

func TestValidateBootstrapOptions(t *testing.T) {
	optionTests := []struct {
		name     string
		baseRepo string
		errMsg   string
	}{
		{"invalid repo", "test", "repo must be org/repo"},
		{"valid repo", "test/repo", ""},
	}

	for _, tt := range optionTests {
		o := BootstrapOptions{quayUsername: "testing", baseRepo: tt.baseRepo, prefix: "test"}

		err := o.Validate()

		if err != nil && tt.errMsg == "" {
			t.Errorf("Validate() %#v got an unexpected error: %s", tt.name, err)
			continue
		}

		if !matchError(t, tt.errMsg, err) {
			t.Errorf("Validate() %#v failed to match error: got %s, want %s", tt.name, err, tt.errMsg)
		}
	}
}

// TODO: set up for complete BootstrapOptions instead of just prefix.
func TestCompleteBootstrapOptions(t *testing.T) {
	completeTests := []struct {
		name       string
		prefix     string
		wantPrefix string
	}{
		{"no prefix", "", ""},
		{"prefix with hyphen", "test-", "test-"},
		{"prefix without hyphen", "test", "test-"},
	}

	for _, tt := range completeTests {
		o := BootstrapOptions{prefix: tt.prefix}

		err := o.Complete("test", &cobra.Command{}, []string{"test", "test/repo"})

		if err != nil {
			t.Fatalf("Complete() %#v failed: ", err)
		}

		if o.prefix != tt.wantPrefix {
			t.Errorf("Complete() %#v prefix: got %s, want %s", tt.name, o.prefix, tt.wantPrefix)
		}
	}
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

func TestBootstrapCommandWithMissingParams(t *testing.T) {
	cmdTests := []struct {
		args    []string
		wantErr string
	}{
		{[]string{}, `Required flag(s) "base-repository", "dockerconfigjson", "github-token", "quay-username" have/has not been set`},

		{[]string{"--base-repository testing/testing --dockerconfigjson testing/test.json --github-token sdfghdfghj --quay-username"}, `Required flag(s) , "quay-username" have/has not been set`},
	}

	for _, tt := range cmdTests {
		cmd := NewCmdBootstrap(BootstrapRecommendedCommandName, util.GetFullName("odo", BootstrapRecommendedCommandName))
		cmd, _, err := executeCommand(cmd)

		if err.Error() != tt.wantErr {
			t.Errorf("got %s, want %s", err, tt.wantErr)
		}
	}
}

func executeCommand(cmd *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs(args)
	c, err = cmd.ExecuteC()
	return c, buf.String(), err
}
