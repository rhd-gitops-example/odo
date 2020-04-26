package config

import (
	"fmt"
	"log"

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
	n, ok := envNames[env.Name]
	if !ok {
		envNames[env.Name] = 1

	} else {
		envNames[env.Name] = n + 1
	}
	if envNames[env.Name] > 1 {
		vv.errs = append(vv.errs, fmt.Errorf("%s environment is more than once", env.Name))
		log.Println("error")
	}
	return nil
}

func (vv *validateVisitor) Application(env *Environment, app *Application) error {
	if err := validName(app.Name); err != nil {
		vv.errs = append(vv.errs, err)
	}
	n, ok := appNames[env.Name+"+"+app.Name]
	if !ok {
		envNames[env.Name+"+"+app.Name] = 1

	} else {
		envNames[env.Name+"+"+app.Name] = n + 1
	}
	if envNames[env.Name+"+"+app.Name] > 1 {
		vv.errs = append(vv.errs, fmt.Errorf("%s app is more than once in environment %s", app.Name, env.Name))
		log.Println("error+1")
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

	n, ok := serviceNames[env.Name+"+"+app.Name+"+"+svc.Name]
	if !ok {
		serviceNames[env.Name+"+"+app.Name+"+"+svc.Name] = 1

	} else {
		serviceNames[env.Name+"+"+app.Name+"+"+svc.Name] = n + 1
	}
	if serviceNames[env.Name+"+"+app.Name+"+"+svc.Name] > 1 {
		vv.errs = append(vv.errs, fmt.Errorf("%s service in %s app is more than once in environment %s", svc.Name, app.Name, env.Name))
		log.Println("error+2")
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
