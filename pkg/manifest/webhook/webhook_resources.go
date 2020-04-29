package webhook

import (
	routeclientset "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/openshift/odo/pkg/manifest/clientconfig"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type webhookResources struct {
	routeClient routeclientset.RouteV1Interface
	kubeClient  kubernetes.Interface
}

// NewResources create new webhook resources
func NewResources() (*webhookResources, error) {
	config, err := clientconfig.GetRESTConfig()
	if err != nil {
		return nil, err
	}
	routeClient, err := routeclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &webhookResources{routeClient: routeClient,
		kubeClient: kubeClient}, nil
}

func (r *webhookResources) getWebhookSecret(ns, secetName, key string) (string, error) {

	secret, err := r.kubeClient.CoreV1().Secrets(ns).Get(secetName, metav1.GetOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "unable to get the secret %s", secret)
	}
	return string(secret.Data[key]), nil
}

func (r *webhookResources) getListenerAddress(ns, routeName string) (string, error) {
	route, err := r.routeClient.Routes(ns).Get(routeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return route.Spec.Host, nil
}
