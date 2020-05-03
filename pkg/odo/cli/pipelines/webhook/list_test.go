package webhook

import (
	"fmt"
	"testing"
)

func TestMissingRequiredFlagsForList(t *testing.T) {

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
			_, _, err := executeCommand(NewCmdList("webhook", "odo pipelines webhook create"), tt.flags...)

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

func TestValidateForList(t *testing.T) {

	testcases := []struct {
		options *listOptions
		errMsg  string
	}{
		{
			&listOptions{
				options{isCICD: true, serviceName: "foo"},
			},
			"Only one of --cicd or --service can be specified",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: ""},
			},
			"One of --cicd or --service must be specified",
		},
		{
			&listOptions{
				options{isCICD: true, serviceName: ""},
			},
			"",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: "foo/bar/gau"},
			},
			"",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: "foo"},
			},
			"Fully qualified service name must be in format <environment name>/<application name>/<service name>",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: "/foo"},
			},
			"Fully qualified service name must be in format <environment name>/<application name>/<service name>",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: "foo/bar/bar/gau"},
			},
			"Fully qualified service name must be in format <environment name>/<application name>/<service name>",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: "/bar/bar"},
			},
			"Fully qualified service name must be in format <environment name>/<application name>/<service name>",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: "bar/foo"},
			},
			"Fully qualified service name must be in format <environment name>/<application name>/<service name>",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: "bar/foo/gau/"},
			},
			"Fully qualified service name must be in format <environment name>/<application name>/<service name>",
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
