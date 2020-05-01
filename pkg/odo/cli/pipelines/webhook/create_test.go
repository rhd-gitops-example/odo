package webhook

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/spf13/cobra"
)

type keyValuePair struct {
	key   string
	value string
}

func TestMissingRequiredFlagsForCreate(t *testing.T) {
	cmdTests := []struct {
		flags   []keyValuePair
		wantErr string
	}{
		{[]keyValuePair{flag("cicd", "true")},
			"Required flag(s) \"access-token\" have/has not been set",
		},
		{[]keyValuePair{flag("access-token", "mytoken"), flag("cicd", "true")},
			"",
		},
	}
	for i, tt := range cmdTests {
		t.Run(fmt.Sprintf("Test %d", i), func(rt *testing.T) {
			_, _, err := executeCommand(NewCmdCreate("webhook", "odo pipelines webhook create"), tt.flags...)

			if err != nil {
				if err.Error() != tt.wantErr {
					rt.Errorf("got %s, want %s", err, tt.wantErr)
				}
			} else {
				if tt.wantErr != "" {
					rt.Errorf("got %s, want %s", "", tt.wantErr)
				}
			}
		})
	}
}

func TestValidateForCreate(t *testing.T) {
	optionTests := []struct {
		options *createOptions
		errMsg  string
	}{
		{
			&createOptions{isCICD: true, serviceName: "abc"},
			"Only one of --cicd or --service-name can be specified",
		},
		{
			&createOptions{isCICD: false, serviceName: ""},
			"One of --cicd or --service-name must be specified",
		},
		{
			&createOptions{isCICD: true, serviceName: ""},
			"",
		},
		{
			&createOptions{isCICD: false, serviceName: "bb"},
			"",
		},
	}

	for i, tt := range optionTests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {

			err := tt.options.Validate()

			if err != nil && tt.errMsg == "" {
				t.Errorf("Validate() got an unexpected error: %s", err)
			} else {
				if !matchError(t, tt.errMsg, err) {
					t.Errorf("Validate() failed to match error: got %s, want %s", err, tt.errMsg)
				}
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

func flag(k, v string) keyValuePair {
	return keyValuePair{
		key:   k,
		value: v,
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
