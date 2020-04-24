package config

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	// "knative.dev/pkg/apis"
)

const (
	DNS1035Error = "[a DNS-1035 label must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character (e.g. 'my-name',  or 'abc-123', regex used for validation is '[a-z]([-a-z0-9]*[a-z0-9])?')]"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		desc string
		file string
		want []error
	}{
		{
			"Misssing required fields error",
			"testdata/name_error.yaml",
			[]error{
				validName(""),
				validName("app-1$."),
				validName("develo.pment"),
				validName("test)cicd"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(rt *testing.T) {
			manifest, err := ParseFile(test.file)
			if err != nil {
				rt.Fatalf("failed to parse file:%#v", err)
			}
			got := manifest.Validate()
			err = matchErrors(rt, got, test.want)
			if err != nil {
				rt.Fatal(err)
			}
		})
	}
}

func matchErrors(t *testing.T, got []error, want []error) error {
	t.Helper()
	if len(got) != len(want) {
		return fmt.Errorf("did not match error, got %v want %v", got, want)
	}
	for i := 0; i < len(got); i++ {
		if diff := cmp.Diff(got[i].Error(), want[i].Error()); diff != "" {
			return fmt.Errorf("did not match error, got %v want %v", got[i], want[i])
		}
	}
	return nil
}
