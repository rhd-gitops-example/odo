package webhook

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/odo/pkg/pipelines/config"
)

func TestBuildURL(t *testing.T) {
	testcases := []struct {
		host     string
		port     string
		insecure bool
		want     string
	}{
		{
			host:     "test.example.com",
			port:     "80",
			insecure: true,
			want:     "http://test.example.com",
		},
		{
			host:     "test.example.com",
			port:     "8080",
			insecure: true,
			want:     "http://test.example.com:8080",
		},
		{
			host:     "test.example.com",
			port:     "443",
			insecure: false,
			want:     "https://test.example.com",
		},
		{
			host:     "test.example.com",
			port:     "8443",
			insecure: false,
			want:     "https://test.example.com:8443",
		},
	}
	for i, tt := range testcases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			got := buildURL(tt.host, tt.port, tt.insecure)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("result mismatch got\n%s", diff)
			}
		})
	}
}
func TestGetGitRepoURL(t *testing.T) {
	testcases := []struct {
		manifest *config.Manifest
		isCICD   bool
		names    []string
		want     string
	}{
		{
			manifest: &config.Manifest{
				GitOpsURL: "https://github.com/foo/bar.git",
			},
			isCICD: true,
			names:  []string{"", "", ""},
			want:   "https://github.com/foo/bar.git",
		},
		{
			manifest: &config.Manifest{},
			want:     "",
		},
		{
			manifest: &config.Manifest{
				GitOpsURL: "https://github.com/foo/bar.git",
				Environments: []*config.Environment{
					{
						Name: "notmyenv",
						Apps: []*config.Application{
							{
								Name: "notmyapp",
								Services: []*config.Service{
									{
										Name:      "notmyservice",
										SourceURL: "https://not/mine",
									},
								},
							},
						},
					},
					{
						Name: "myenv",
						Apps: []*config.Application{
							{
								Name: "notmyapp",
								Services: []*config.Service{
									{
										Name:      "notmyservice",
										SourceURL: "https://not/mine",
									},
								},
							},
							{
								Name: "myapp",
								Services: []*config.Service{
									{
										Name:      "notmyservice",
										SourceURL: "https://not/mine",
									},
									{
										Name:      "myservice",
										SourceURL: "https://github.com/foo2/bar.git",
									},
									{
										Name:      "notmyserviceagain",
										SourceURL: "https://not/mine",
									},
								},
							},
							{
								Name: "notmyapp2",
								Services: []*config.Service{
									{
										Name:      "notmyservice",
										SourceURL: "https://not/mine",
									},
								},
							},
						},
					},
				},
			},
			isCICD: false,
			names:  []string{"myenv", "myapp", "myservice"},
			want:   "https://github.com/foo2/bar.git",
		},
	}
	for i, tt := range testcases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			got := getRepoURL(tt.manifest, tt.isCICD, tt.names)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("result mismatch got\n%s", diff)
			}
		})
	}
}

func TestGetCICDEnvName(t *testing.T) {
	testcases := []struct {
		manifest *config.Manifest
		want     string
	}{
		{
			manifest: &config.Manifest{},
			want:     "",
		},
		{
			manifest: &config.Manifest{
				GitOpsURL: "https://github.com/foo/bar.git",
				Environments: []*config.Environment{
					{
						Name:   "notmyenv",
						IsCICD: false,
						Apps: []*config.Application{
							{
								Name: "notmyapp",
								Services: []*config.Service{
									{
										Name:      "notmyservice",
										SourceURL: "https://not/mine",
									},
								},
							},
						},
					},
					{
						Name: "myenv",
						Apps: []*config.Application{
							{
								Name: "notmyapp",
								Services: []*config.Service{
									{
										Name:      "notmyservice",
										SourceURL: "https://not/mine",
									},
								},
							},
							{
								Name: "myapp",
								Services: []*config.Service{
									{
										Name:      "notmyservice",
										SourceURL: "https://not/mine",
									},
									{
										Name:      "myservice",
										SourceURL: "https://github.com/foo2/bar.git",
									},
									{
										Name:      "notmyserviceagain",
										SourceURL: "https://not/mine",
									},
								},
							},
							{
								Name: "notmyapp2",
								Services: []*config.Service{
									{
										Name:      "notmyservice",
										SourceURL: "https://not/mine",
									},
								},
							},
						},
					},
					{
						Name:   "notmyenv2",
						IsCICD: true,
						Apps: []*config.Application{
							{
								Name: "notmyapp",
								Services: []*config.Service{
									{
										Name:      "notmyservice",
										SourceURL: "https://not/mine",
									},
								},
							},
						},
					},
				},
			},
			want: "notmyenv2",
		},
	}
	for i, tt := range testcases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			got := getCICDNamespace(tt.manifest)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("result mismatch got\n%s", diff)
			}
		})
	}
}
