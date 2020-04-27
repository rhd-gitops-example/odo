package config

import (
	"fmt"

	"path/filepath"

	"github.com/mkmik/multierror"
	"k8s.io/apimachinery/pkg/api/validation"
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
	if err := validateName(env.Name, PathForEnvironment(env)); err != nil {
		vv.errs = append(vv.errs, err)
	}
	if err := validatePipelines(env.Pipelines, PathForEnvironment(env)); err != nil {
		vv.errs = append(vv.errs, err...)
	}
	error := validateEnvironments(env.Name)
	if error != nil {
		vv.errs = append(vv.errs, error)
	}
	return nil
}

func (vv *validateVisitor) Application(env *Environment, app *Application) error {
	if err := validateName(app.Name, PathForApplication(env, app)); err != nil {
		vv.errs = append(vv.errs, err)
	}
	errors := validateApplications(env.Name, app.Name)
	if errors != nil {
		vv.errs = append(vv.errs, errors)
	}
	return nil
}

func (vv *validateVisitor) Service(env *Environment, app *Application, svc *Service) error {
	if err := validateName(svc.Name, PathForService(env, svc)); err != nil {
		vv.errs = append(vv.errs, err)
	}
	if err := validateWebhook(svc.Webhook, PathForService(env, svc)); err != nil {
		vv.errs = append(vv.errs, err...)
	}
	if err := validatePipelines(svc.Pipelines, PathForService(env, svc)); err != nil {
		vv.errs = append(vv.errs, err...)
	}

	errors := validateService(app.Name, svc.Name)
	if errors != nil {
		vv.errs = append(vv.errs, errors)
	}
	return nil
}

func validateWebhook(hook *Webhook, path string) []error {
	errs := []error{}
	if hook == nil {
		return nil
	}
	if hook.Secret == nil {
		return list(notFoundError("secret", path))
	}
	if err := validateName(hook.Secret.Name, filepath.Join(path, "secret")); err != nil {
		errs = append(errs, err)
	}
	if err := validateName(hook.Secret.Namespace, filepath.Join(path, "secret")); err != nil {
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
		return list(notFoundError("pipelines", path))
	}
	if err := validateName(pipelines.Integration.Template, filepath.Join(path, "pipelines")); err != nil {
		errs = append(errs, err)
	}
	if err := validateName(pipelines.Integration.Binding, filepath.Join(path, "pipelines")); err != nil {
		errs = append(errs, err)
	}
	return errs
}

func validateName(name, path string) error {
	err := validation.NameIsDNS1035Label(name, true)
	if len(err) > 0 {
		return fmt.Errorf("%q is not a valid name at %v: \n%v", name, path, err)
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

func notFoundError(field string, at string) error {
	return fmt.Errorf("%v not found at %v", field, at)
}

func list(errs ...error) []error {
	return errs
}
