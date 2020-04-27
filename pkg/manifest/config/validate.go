package config

import (
	"fmt"

	"strings"

	"github.com/mkmik/multierror"
	"k8s.io/apimachinery/pkg/api/validation"
	"knative.dev/pkg/apis"
)

var envNames = map[string]int{}
var appNames = map[string]int{}
var serviceNames = map[string]int{}

type validateVisitor struct {
	errs []error
}

func (m *Manifest) Validate() error {
	vv := &validateVisitor{errs: []error{}}
	m.Walk(vv)
	return multierror.Join(vv.errs)
}

func (vv *validateVisitor) Environment(env *Environment) error {
	envPath := yamlPath(PathForEnvironment(env))
	if err := validateName(env.Name, envPath); err != nil {
		vv.errs = append(vv.errs, err)
	}
	if err := validatePipelines(env.Pipelines, envPath); err != nil {
		vv.errs = append(vv.errs, err...)
	}

	// error := validateEnvironments(env.Name)
	// if error != nil {
	// 	vv.errs = append(vv.errs, error)
	// }
	return nil
}

func (vv *validateVisitor) Application(env *Environment, app *Application) error {
	appPath := yamlPath(PathForApplication(env, app))
	if err := validateName(app.Name, appPath); err != nil {
		vv.errs = append(vv.errs, err)
	}

	// errors := validateApplications(env.Name, app.Name)
	// if errors != nil {
	// 	vv.errs = append(vv.errs, errors)
	// }
	return nil
}

func (vv *validateVisitor) Service(env *Environment, app *Application, svc *Service) error {
	svcPath := yamlPath(PathForService(env, svc))
	if err := validateName(svc.Name, svcPath); err != nil {
		vv.errs = append(vv.errs, err)
	}
	// errors := validateService(app.Name, svc.Name)
	// if errors != nil {
	// 	vv.errs = append(vv.errs, errors)

	if err := validateWebhook(svc.Webhook, svcPath); err != nil {
		vv.errs = append(vv.errs, err...)
	}
	if err := validatePipelines(svc.Pipelines, svcPath); err != nil {
		vv.errs = append(vv.errs, err...)
	}
	return nil
}

func validateWebhook(hook *Webhook, path string) []error {
	errs := []error{}
	if hook == nil {
		return nil
	}
	if hook.Secret == nil {
		return list(apis.ErrMissingField(yamlJoin(path, "webhook", "secret")))
	}
	if err := validateName(hook.Secret.Name, yamlJoin(path, "webhook", "secret", "name")); err != nil {
		errs = append(errs, err)
	}
	if err := validateName(hook.Secret.Namespace, yamlJoin(path, "webhook", "secret", "namespace")); err != nil {
		errs = append(errs, err)
	}
	return errs
}

func validatePipelines(pipelines *Pipelines, path string) []error {
	errs := []error{}
	if pipelines == nil {
		return nil
	}
	if pipelines.Integration == nil {
		return list(apis.ErrMissingField(yamlJoin(path, "pipelines", "integration")))
	}
	if err := validateName(pipelines.Integration.Template, yamlJoin(path, "pipelines", "integration", "template")); err != nil {
		errs = append(errs, err)
	}
	if err := validateName(pipelines.Integration.Binding, yamlJoin(path, "pipelines", "integration", "binding")); err != nil {
		errs = append(errs, err)
	}
	return errs
}

func validateName(name, path string) *apis.FieldError {
	err := validation.NameIsDNS1035Label(name, true)
	if len(err) > 0 {
		return invalidNameError(name, err[0], []string{path})
	}
	return nil
}

func validateEnvironments(envName string) error {
	n, ok := envNames[envName]
	if !ok {
		envNames[envName] = 1

	} else {
		envNames[envName] = n + 1
	}
	if envNames[envName] > 1 {
		return fmt.Errorf("%s environment is more than once", envName)
	}
	return nil
}

func validateApplications(envName, appName string) error {
	n, ok := appNames[envName+"+"+appName]
	if !ok {
		appNames[envName+"+"+appName] = 1

	} else {
		appNames[envName+"+"+appName] = n + 1
	}
	if appNames[envName+"+"+appName] > 1 {
		return fmt.Errorf("%s app is more than once in environment %s", appName, envName)

	}
	return nil
}

func validateService(appName, svcName string) error {
	n, ok := serviceNames[appName+"+"+svcName]
	if !ok {
		serviceNames[appName+"+"+svcName] = 1

	} else {
		serviceNames[appName+"+"+svcName] = n + 1
	}
	if serviceNames[appName+"+"+svcName] > 1 {
		return fmt.Errorf("%s service in %s app is more than once", svcName, appName)
	}
	return nil

}

func yamlPath(path string) string {
	return strings.ReplaceAll(path, "/", ".")
}

func yamlJoin(a string, b ...string) string {
	for _, s := range b {
		a = a + "." + s
	}
	return a
}

func list(errs ...error) []error {
	return errs
}

func invalidNameError(name, details string, paths []string) *apis.FieldError {
	return &apis.FieldError{
		Message: fmt.Sprintf("invalid name %q", name),
		Details: details,
		Paths:   paths,
	}
}
