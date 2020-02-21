package pipelines

import (
	"bytes"
	"testing"

	trial "github.com/openshift/odo/pkg/odo/cli/pipelines"
	"github.com/spf13/cobra"
)

func TestBootstrapCommandWithMissingParams(t *testing.T) {
	cmdTests := []struct {
		args    []string
		wantErr string
	}{
		{[]string{"quay-username", "example", "github-token", "abc123", "dockerconfigjson", "~/"}, `Required flag(s) "base-repository" have/has not been set`},
		{[]string{"quay-username", "example", "github-token", "abc123", "base-repository", "example/repo"}, `Required flag(s) "dockerconfigjson" have/has not been set`},
		{[]string{"quay-username", "example", "dockerconfigjson", "~/", "base-repository", "example/repo"}, `Required flag(s) "github-token" have/has not been set`},
		{[]string{"github-token", "abc123", "dockerconfigjson", "~/", "base-repository", "example/repo"}, `Required flag(s) "quay-username" have/has not been set`},
	}
	for _, tt := range cmdTests {
		_, _, err := executeCommand(trial.NewCmdBootstrap("bootstrap", "odo pipelines bootstrap"), tt.args...)
		if err.Error() != tt.wantErr {
			t.Errorf("got %s, want %s", err, tt.wantErr)
		}
	}
}

func executeCommand(cmd *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.Flags().Set(args[0], args[1])
	cmd.Flags().Set(args[2], args[3])
	cmd.Flags().Set(args[4], args[5])
	c, err = cmd.ExecuteC()
	return c, buf.String(), err
}
