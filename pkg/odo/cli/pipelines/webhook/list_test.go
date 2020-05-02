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
		{[]keyValuePair{flag("access-token", "mytoken"), flag("cicd", "true")},
			"",
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
			"Only one of --cicd or --service-name can be specified",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: ""},
			},
			"One of --cicd or --service-name must be specified",
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
			"Fully qualifed service-name must be in format <application name>/<service name>",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: "/foo"},
			},
			"Fully qualifed service-name must be in format <application name>/<service name>",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: "foo/bar/bar/gau"},
			},
			"Fully qualifed service-name must be in format <application name>/<service name>",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: "/bar/bar"},
			},
			"Fully qualifed service-name must be in format <application name>/<service name>",
		},
		{
			&listOptions{
				options{isCICD: false, serviceName: "bar/foo"},
			},
			"Fully qualifed service-name must be in format <application name>/<service name>",
		},
		{
			&listOptions{
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
