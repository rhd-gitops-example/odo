package webhook

import (
	routeclientset "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/openshift/odo/pkg/manifest/clientconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type webhookResources struct {
	routeClient routeclientset.RouteV1Interface
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
	return &webhookResources{routeClient: routeClient}, nil
}

func (r *webhookResources) getWebhookSecret(ns, secetName string) (string, error) {
	return "", nil
}

func (r *webhookResources) getListenerAddress(ns, routeName string) (string, error) {
	route, err := r.routeClient.Routes(ns).Get(routeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return route.Spec.Host, nil
}
