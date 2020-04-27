package config

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/validation"
)

var envNames = map[string]int{}
var appNames = map[string]int{}
var serviceNames = map[string]int{}

type validateVisitor struct {
	errs []error
}

func (m *Manifest) Validate() []error {
	vv := &validateVisitor{errs: []error{}}
	if err := m.Walk(vv); err != nil {
		return list(err)
	}
	return vv.errs
}

func list(errs ...error) []error {
	return errs
}

func (vv *validateVisitor) Environment(env *Environment) error {
	if err := validName(env.Name); err != nil {
		vv.errs = append(vv.errs, err)
	}

	error := validateEnvironments(env.Name)
	if error != nil {
		vv.errs = append(vv.errs, error)
	}

	return nil
}

func (vv *validateVisitor) Application(env *Environment, app *Application) error {
	if err := validName(app.Name); err != nil {
		vv.errs = append(vv.errs, err)
	}
	errors := validateApplications(env.Name, app.Name)
	if errors != nil {
		vv.errs = append(vv.errs, errors)
	}

	return nil
}

func (vv *validateVisitor) Service(env *Environment, app *Application, svc *Service) error {
	if err := validName(svc.Name); err != nil {
		vv.errs = append(vv.errs, err)
	}
	if err := validateWebhook(svc); err != nil {
		return err
	}

	if err := validatePipelines(svc); err != nil {
		return err
	}

	errors := validateService(app.Name, svc.Name)
	if errors != nil {
		vv.errs = append(vv.errs, errors)
	}
	return nil
}

func validateWebhook(svc *Service) error {
	if svc.Webhook == nil {
		return nil
	}
	if svc.Webhook.Secret == nil {
		return notFoundError([]string{"secret"}, svc.Name)
	}
	if err := validName(svc.Webhook.Secret.Name); err != nil {
		return err
	}
	if err := validName(svc.Webhook.Secret.Namespace); err != nil {
		return err
	}

	return nil
}

func validatePipelines(svc *Service) error {
	if svc.Pipelines == nil {
		return nil
	}
	if svc.Pipelines.Integration == nil {
		return notFoundError([]string{"templates", "bindings"}, svc.Name)
	}
	if err := validName(svc.Pipelines.Integration.Template); err != nil {
		return err
	}
	if err := validName(svc.Pipelines.Integration.Binding); err != nil {
		return err
	}
	return nil
}

func validName(name string) error {
	err := validation.NameIsDNS1035Label(name, true)
	if len(err) > 0 {
		return fmt.Errorf("%q is not a valid name: %v", name, err)
	}
	return nil
}

func notFoundError(items []string, at string) error {
	return fmt.Errorf("%v not found at %v", items, at)
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
