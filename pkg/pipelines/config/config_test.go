package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestManifestWalk(t *testing.T) {
	m := &Manifest{
		Config: &Config{
			CICD: &Cicd{
				Namespace: "cicd",
			},
			Argo: &Argo{
				Namespace: "argocd",
			},
		},
		Environments: []*Environment{
			{
				Name: "development",
				Services: []*Service{
					{Name: "app-1-service-http"},
					{Name: "app-1-service-test"},
					{Name: "app-2-service"},
				},
				Apps: []*Application{
					{
						Name: "my-app-1",
						ServiceRefs: []string{
							"app-1-service-http",
							"app-1-service-test",
						},
					},
					{
						Name: "my-app-2",
						ServiceRefs: []string{
							"app-2-service",
						},
					},
				},
			},
			{
				Name: "staging",
				Services: []*Service{
					{Name: "app-1-service-user"},
				},
				Apps: []*Application{
					{Name: "my-app-1",
						ServiceRefs: []string{
							"app-1-service-user",
						},
					},
				},
			},
		},
	}

	v := &testVisitor{paths: []string{}}
	err := m.Walk(v)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(v.paths)

	want := []string{
		"development/app-1-service-http",
		"development/app-1-service-test",
		"development/app-2-service",
		"development/my-app-1",
		"development/my-app-2",
		"envs/development",
		"envs/staging",
		"staging/app-1-service-user",
		"staging/my-app-1",
	}

	if diff := cmp.Diff(want, v.paths); diff != "" {
		t.Fatalf("tree files: %s", diff)
	}
}

func TestEnviromentSorting(t *testing.T) {
	envNames := func(envs []*Environment) []string {
		n := make([]string, len(envs))
		for i, v := range envs {
			n[i] = v.Name
		}
		return n
	}
	envTests := []struct {
		names []testEnv
		want  []string
	}{
		{[]testEnv{{"prod", false, false}, {"staging", false, false}, {"dev", false, false}}, []string{"dev", "prod", "staging"}},
	}

	for _, tt := range envTests {
		envs := makeEnvs(tt.names)
		sort.Sort(ByName(envs))
		if diff := cmp.Diff(tt.want, envNames(envs)); diff != "" {
			t.Errorf("sort(%#v): %s", envs, diff)
		}
	}
}

func TestFindCICDEnviroment(t *testing.T) {
	envTests := []struct {
		manifest *Manifest
		want     *Cicd
		err      string
	}{
		{&Manifest{
			Config: &Config{
				CICD: &Cicd{
					Namespace: "cicd",
				},
			},
		}, &Cicd{
			Namespace: "cicd",
		}, ""},
		{&Manifest{
			Config: &Config{
				Argo: &Argo{
					Namespace: "argocd",
				},
			},
		}, nil, ""},
	}

	for i, tt := range envTests {
		t.Run(fmt.Sprintf("test %d", i), func(rt *testing.T) {
			m := tt.manifest
			_, err := m.GetCICD()
			if !matchErrorString(t, tt.err, err) {
				rt.Errorf("did not match error, got %s, want %s", err, tt.err)
				return
			}
		})
	}
}

func TestGetEnvironment(t *testing.T) {
	m := &Manifest{Environments: makeEnvs([]testEnv{{name: "prod"}, {name: "testing"}})}
	env := m.GetEnvironment("prod")
	if env.Name != "prod" {
		t.Fatalf("got the wrong environment back: %#v", env)
	}

	unknown := m.GetEnvironment("unknown")
	if unknown != nil {
		t.Fatalf("found an unknown env: %#v", unknown)
	}
}

func makeEnvs(ns []testEnv) []*Environment {
	n := make([]*Environment, len(ns))
	for i, v := range ns {
		n[i] = &Environment{Name: v.name}
	}
	return n

}

type testEnv struct {
	name   string
	cicd   bool
	argocd bool
}

type testVisitor struct {
	pipelineServices []string
	paths            []string
}

func (v *testVisitor) Service(env *Environment, svc *Service) error {
	v.paths = append(v.paths, filepath.Join(env.Name, svc.Name))
	v.pipelineServices = append(v.pipelineServices, filepath.Join("cicd", env.Name, svc.Name))
	return nil
}

func (v *testVisitor) Application(env *Environment, app *Application) error {
	v.paths = append(v.paths, filepath.Join(env.Name, app.Name))
	return nil
}

func (v *testVisitor) Environment(env *Environment) error {
	if env.Name == "cicd" {
		v.paths = append(v.paths, v.pipelineServices...)
	}
	v.paths = append(v.paths, filepath.Join("envs", env.Name))
	return nil
}

// MatchErrorString takes a string and matches on the error and returns true if
// the
// string matches the error.
//
// This is useful in table tests.
//
// If the string can't be compiled as an regexp, then this will fail with a
// Fatal error.
func matchErrorString(t *testing.T, s string, e error) bool {
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
