package pipelines

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidateImageRepo(t *testing.T) {
	tests := []struct {
		description string
		URL         string
		validURL    string
		isValid     bool
	}{
		{
			"Valid image regsitry URL",
			"quay.io/sample-user/sample-repo",
			"quay.io/sample-user/sample-repo",
			true,
		},
		{
			"Invalid image registry URL with missing repo name",
			"quay.io/sample-user",
			"quay.io/sample-user/sample-repo",
			false,
		},
		{
			"Invalid image registry URL with whitespaces",
			"quay.io/sample-user/ ",
			"quay.io/sample-user/sample-repo",
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			isValid, _ := validateImageRepo(test.URL)
			if diff := cmp.Diff(isValid, test.isValid); diff != "" {
				t.Errorf("validateImageRepo() failed:\n%s", diff)
			}
		})
	}
}

func TestCheckInternalRegistry(t *testing.T) {
	tests := []struct {
		Description string
		URL         string
		Result      bool
	}{
		{
			"Valid internal registry URL",
			"image-registry.openshift-image-registry.svc:5000/project/app",
			true,
		},
		{
			"Invalid internal registry URL",
			"quay.io/project/app",
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			usingInternalRegistry := checkInternalRegistry(test.URL)
			if usingInternalRegistry != test.Result {
				t.Errorf("validateImageRepo() failed:%s, expected %v but got %v", test.Description, test.Result, usingInternalRegistry)
			}
		})
	}
}
