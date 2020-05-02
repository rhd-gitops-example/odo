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
	testcases := []struct {
		flags   []keyValuePair
		wantErr string
	}{
		{[]keyValuePair{flag("cicd", "true")},
			"Required flag(s) \"access-token\" have/has not been set",
		},
	}
	for i, tt := range testcases {
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
	testcases := []struct {
		options *createOptions
		errMsg  string
	}{
		{
			&createOptions{
				options{isCICD: true, serviceName: "foo"},
			},
			"Only one of --cicd or --service-name can be specified",
		},
		{
			&createOptions{
				options{isCICD: false, serviceName: ""},
			},
			"One of --cicd or --service-name must be specified",
		},
		{
			&createOptions{
				options{isCICD: true, serviceName: ""},
			},
			"",
		},
		{
			&createOptions{
				options{isCICD: false, serviceName: "foo/bar/gau"},
			},
			"",
		},
		{
			&createOptions{
				options{isCICD: false, serviceName: "foo"},
			},
			"Fully qualifed service-name must be in format <application name>/<service name>",
		},
		{
			&createOptions{
				options{isCICD: false, serviceName: "/foo"},
			},
			"Fully qualifed service-name must be in format <application name>/<service name>",
		},
		{
			&createOptions{
				options{isCICD: false, serviceName: "foo/bar/bar/gau"},
			},
			"Fully qualifed service-name must be in format <application name>/<service name>",
		},
		{
			&createOptions{
				options{isCICD: false, serviceName: "/bar/bar"},
			},
			"Fully qualifed service-name must be in format <application name>/<service name>",
		},
		{
			&createOptions{
				options{isCICD: false, serviceName: "bar/foo"},
			},
			"Fully qualifed service-name must be in format <application name>/<service name>",
		},
		{
			&createOptions{
				options{isCICD: false, serviceName: "bar/foo/gau/"},
			},
			"Fully qualifed service-name must be in format <application name>/<service name>",
		},
	}

	for i, tt := range testcases {
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
