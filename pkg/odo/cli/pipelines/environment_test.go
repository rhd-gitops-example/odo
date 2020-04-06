package pipelines

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCompleteEnvParameters(t *testing.T) {
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
		o := EnvParameters{prefix: tt.prefix}

		err := o.Complete("test", &cobra.Command{}, []string{"test", "test/repo"})

		if err != nil {
			t.Errorf("Complete() %#v failed: ", err)
		}

		if o.prefix != tt.wantPrefix {
			t.Errorf("Complete() %#v prefix: got %s, want %s", tt.name, o.prefix, tt.wantPrefix)
		}
	}
}

func TestValidateEnvParameters(t *testing.T) {
	optionTests := []struct {
		name    string
		gitRepo string
		errMsg  string
	}{
		{"invalid repo", "test", "repo must be org/repo"},
		{"valid repo", "test/repo", ""},
	}

	for _, tt := range optionTests {
		o := EnvParameters{gitOpsRepo: tt.gitRepo, prefix: "test"}

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

func TestAddCommandWithMissingParams(t *testing.T) {
	cmdTests := []struct {
		desc    string
		flags   []keyValuePair
		wantErr string
	}{
		{"Missing gitops-repo flag",
			[]keyValuePair{flag("output", "~/output"),
				flag("env-name", "test"), flag("skip-checks", "true")},
			`Required flag(s) "gitops-repo" have/has not been set`},
		{"Missing env-name flag",
			[]keyValuePair{flag("gitops-repo", "org/sample"), flag("output", "~/output"),
				flag("skip-checks", "true")},
			`Required flag(s) "env-name" have/has not been set`},
	}
	for _, tt := range cmdTests {
		t.Run(tt.desc, func(t *testing.T) {
			_, _, err := executeCommand(NewCmdEnv("add", "odo pipelines add"), tt.flags...)
			if err.Error() != tt.wantErr {
				t.Errorf("got %s, want %s", err, tt.wantErr)
			}
		})
	}
}
