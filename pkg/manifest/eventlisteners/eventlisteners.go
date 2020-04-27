package eventlisteners

import (
	"fmt"

	"github.com/openshift/odo/pkg/manifest/meta"
	"github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	saName = "pipeline"
)

// Filters for interceptors
const (
	StageCIDryRunFilters = "(header.match('X-GitHub-Event', 'pull_request') && body.action == 'opened' || body.action == 'synchronize') && body.pull_request.head.repo.full_name == '%s'"

	StageCDDeployFilters = "(header.match('X-GitHub-Event', 'push') && body.repository.full_name == '%s') && body.ref.startsWith('refs/heads/master')"

	GitOpsWebhookSecret = "gitops-webhook-secret"

	WebhookSecretKey = "webhook-secret-key"
)

var (
	eventListenerTypeMeta = meta.TypeMeta("EventListener", "tekton.dev/v1alpha1")
)

// Generate will create the required eventlisteners.
func Generate(githubRepo, ns, saName, secretName string) triggersv1.EventListener {
	return triggersv1.EventListener{
		TypeMeta:   eventListenerTypeMeta,
		ObjectMeta: createListenerObjectMeta("cicd-event-listener", ns),
		Spec: triggersv1.EventListenerSpec{
			ServiceAccountName: saName,
			Triggers: []triggersv1.EventListenerTrigger{
				CreateListenerTrigger(
					"ci-dryrun-from-pr",
					StageCIDryRunFilters,
					githubRepo,
					"github-pr-binding",
					"ci-dryrun-from-pr-template",
					secretName,
					ns,
				),
				CreateListenerTrigger(
					"cd-deploy-from-push",
					StageCDDeployFilters,
					githubRepo,
					"github-push-binding",
					"cd-deploy-from-push-template",
					secretName,
					ns,
				),
			},
		},
	}
}

func CreateELFromTriggers(cicdNs string, triggers []triggersv1.EventListenerTrigger) *triggersv1.EventListener {
	return &v1alpha1.EventListener{
		TypeMeta: eventListenerTypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cicd-event-listener",
			Namespace: cicdNs,
		},
		Spec: triggersv1.EventListenerSpec{
			ServiceAccountName: saName,
			Triggers:           triggers,
		},
	}
}

func createEventInterceptor(filter string, repoName string) *triggersv1.EventInterceptor {
	return &triggersv1.EventInterceptor{
		CEL: &triggersv1.CELInterceptor{
			Filter: fmt.Sprintf(filter, repoName),
		},
	}
}

func createGitHubInterceptor(secretName, ns string) *triggersv1.EventInterceptor {
	return &triggersv1.EventInterceptor{
		GitHub: &triggersv1.GitHubInterceptor{
			SecretRef: &triggersv1.SecretRef{
				SecretName: secretName,
				SecretKey:  WebhookSecretKey,
				Namespace:  ns,
			},
		},
	}
}

func CreateListenerTrigger(name, filter, repoName, binding, template, secretName, secretNs string) triggersv1.EventListenerTrigger {
	return triggersv1.EventListenerTrigger{
		Name: name,
		Interceptors: []*triggersv1.EventInterceptor{
			createEventInterceptor(filter, repoName),
			createGitHubInterceptor(secretName, secretNs),
		},
		Bindings: []*triggersv1.EventListenerBinding{
			createListenerBinding(binding),
		},
		Template: createListenerTemplate(template),
	}
}

func createListenerTemplate(name string) triggersv1.EventListenerTemplate {
	return triggersv1.EventListenerTemplate{
		Name: name,
	}
}

func createListenerBinding(name string) *triggersv1.EventListenerBinding {
	return &triggersv1.EventListenerBinding{
		Name: name,
	}
}

func createListenerObjectMeta(name, ns string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}
}