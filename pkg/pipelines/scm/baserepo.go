package scm
import(
	"net/url"
	"strings"

	"github.com/openshift/odo/pkg/pipelines/meta"
	"github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"

)
// NewGitLabRepository returns an instance of GitLabRepository
type processRawURL(rawURL string, func processPath(parsedURL *url.URL) (string, err)){
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	path, err := processPath(parsedURL)
	if err != nil {
		return nil, err
	}
}