package pipelines

import (
	"errors"
	"fmt"
	"os"
	"path"

	corev1 "k8s.io/api/core/v1"
	v1rbac "k8s.io/api/rbac/v1"

	"github.com/mitchellh/go-homedir"
	"github.com/openshift/odo/pkg/pipelines/eventlisteners"
	"github.com/openshift/odo/pkg/pipelines/routes"
	"github.com/openshift/odo/pkg/pipelines/tasks"
	"sigs.k8s.io/yaml"
)

<<<<<<< HEAD
<<<<<<< HEAD
var (
	dockerSecretName = "regcred"
	saName           = "demo-sa"
	roleName         = "tekton-triggers-openshift-demo"
	roleBindingName  = "tekton-triggers-openshift-binding"

	// PolicyRules to be bound to service account
	rules = []v1rbac.PolicyRule{
		v1rbac.PolicyRule{
			APIGroups: []string{"tekton.dev"},
			Resources: []string{"eventlisteners", "triggerbindings", "triggertemplates", "tasks", "taskruns"},
			Verbs:     []string{"get"},
		},
		v1rbac.PolicyRule{
			APIGroups: []string{"tekton.dev"},
			Resources: []string{"pipelineruns", "pipelineresources", "taskruns"},
			Verbs:     []string{"create"},
		},
	}
)

// DefaultTokenFileName is the default token file name
var DefaultTokenFileName string = "~/Downloads/token/token.txt"

// DefaultQuayIOAuthFileName is the default auth file name
var DefaultQuayIOAuthFileName string = "~/Downloads/<username>-auth.json"
=======
var DefaultTokenFilename string = "~/Downloads/token"
var DefaultQuayIOAuthFilename string = "~/Downloads/<username>-auth.json"
>>>>>>> de06981f0... made the requested changes mentioned in the PR

=======
>>>>>>> d8c543602... created 4 required-flags quay-username, github-token, dockerconfigjson,base-repository
// BootstrapOptions is a struct that provides the optional flags
type BootstrapOptions struct {
	Prefix string
}

// Bootstrap is the main driver for getting OpenShift pipelines for GitOps
// configured with a basic configuration.
func Bootstrap(quayUsername, baseRepo, githubToken, quayIOAuthFilename string, o *BootstrapOptions) error {

	// First, check for Tekton.  We proceed only if Tekton is installed
	installed, err := checkTektonInstall()
	if err != nil {
		return fmt.Errorf("failed to run Tekton Pipelines installation check: %w", err)
	}
	if !installed {
		return errors.New("failed due to Tekton Pipelines or Triggers are not installed")
	}

	outputs := make([]interface{}, 0)
<<<<<<< HEAD
	//  Create GitHub Secret
	githubAuth, err := createGithubSecret()
	tokenPath, err := getTokenFileName(o)
=======

<<<<<<< HEAD
	tokenPath, err := homedir.Expand(o.TokenFilename)
>>>>>>> 875ecd79c... made the requested changes mentioned in the PR and removed the unwanted comments
	if err != nil {
		return fmt.Errorf("failed to generate token path to file: %w", err)
	}
	f, err := os.Open(tokenPath)
	if err != nil {
<<<<<<< HEAD
		return fmt.Errorf("failed to open path to TokenFileName: %w", err)
=======
		return fmt.Errorf("failed to expand the path %s: %w", authJSONPath, err)
>>>>>>> 424497f3b... made the required changes mentioned in the PR
	}
	outputs = append(outputs, githubAuth)

	// Create Docker Secret
	dockerSecret, err := createDockerSecret(quayUsername)
=======
	githubAuth, err := createOpaqueSecret("github-auth", githubToken)
>>>>>>> d8c543602... created 4 required-flags quay-username, github-token, dockerconfigjson,base-repository
	if err != nil {
<<<<<<< HEAD
		return err
=======
		return fmt.Errorf("failed to expand the path %s: %w", authJSONPath, err)
>>>>>>> 424497f3b... made the required changes mentioned in the PR
	}
	if o.QuayIOAuthFileName == "~/Downloads/<username>-auth.json" {
		authJSONPath, err := pathToDownloadedFile(quayUsername + "-auth.json")
		if err != nil {
			return fmt.Errorf("failed to generate path to file: %w", err)
		}

		f, err = os.Open(authJSONPath)
		if err != nil {
			return err
		}
	} else {
		f, err = os.Open(o.QuayIOAuthFileName)
		if err != nil {
			return err
		}
	}
	outputs = append(outputs, dockerSecret)

	tasks := tasks.Generate(githubAuth.GetName())
	for _, task := range tasks {
		outputs = append(outputs, task)
	}

	eventListener := eventlisteners.Generate(baseRepo)
	outputs = append(outputs, eventListener)

	route := routes.Generate()
	outputs = append(outputs, route)

	//  Create Service Account, Role, Role Bindings, and ClusterRole Bindings
	sa := createServiceAccount(saName, dockerSecretName)
	outputs = append(outputs, sa)
	role := createRole(roleName, rules)
	outputs = append(outputs, role)
	outputs = append(outputs, createRoleBinding(roleBindingName, &sa, role.Kind, role.Name))
	outputs = append(outputs, createRoleBinding("edit-clusterrole-binding", &sa, "ClusterRole", "edit"))

	// Marshall
	for _, r := range outputs {
		data, err := yaml.Marshal(r)
		if err != nil {
			return err
		}
		fmt.Printf("%s---\n", data)
	}

	return nil
}

// createGithubSecret creates Github secret
func createGithubSecret() (*corev1.Secret, error) {
	tokenPath, err := pathToDownloadedFile("token")
	if err != nil {
		return nil, fmt.Errorf("failed to generate path to file: %w", err)
	}
	f, err := os.Open(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file %s due to %w", tokenPath, err)
	}
	outputs = append(outputs, githubAuth)

	authJSONPath, err := homedir.Expand(quayIOAuthFilename)
	if err != nil {
		return fmt.Errorf("failed to generate path to file: %w", err)
	}

	f, err := os.Open(authJSONPath)
	if err != nil {
		return fmt.Errorf("failed to open path authJSON : %w", err)
	}
	defer f.Close()

	githubAuth, err := createOpaqueSecret("github-auth", f)
	if err != nil {
		return nil, err
	}

	return githubAuth, nil
}

// createDockerSecret creates Docker secret
func createDockerSecret(quayUsername string) (*corev1.Secret, error) {
	authJSONPath, err := pathToDownloadedFile(quayUsername + "-auth.json")
	if err != nil {
		return nil, fmt.Errorf("failed to generate path to file: %w", err)
	}

	f, err := os.Open(authJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read docker file '%s' due to %w", authJSONPath, err)
	}
	defer f.Close()

	dockerSecret, err := createDockerConfigSecret(dockerSecretName, f)
	if err != nil {
		return nil, err
	}
	return dockerSecret, nil
}

func pathToDownloadedFile(fname string) (string, error) {
	return homedir.Expand(path.Join("~/Downloads/", fname))
}

// create and invoke a Tekton Checker
func checkTektonInstall() (bool, error) {
	tektonChecker, err := newTektonChecker()
	if err != nil {
		return false, err
	}
	return tektonChecker.checkInstall()
}
