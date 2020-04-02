package pipelines
 
import (
   "errors"
   "fmt"
   "path/filepath"
)
 
// InitParameters is a struct that provides flags for initialise command
type AddParameters struct {
   GitopsRepo          string
   GitopsWebhookSecret string
   Output              string
   Prefix              string
   AppGitRepo          string
   AppWebhookSecret    string
   AppImageRepo        string
   EnvName             string
   DockerCfgJson       string
   SkipChecks          bool
}
 
const (
   appsDir  = "apps"
   overlays = "overlays"
)
 
// Init function will initialise the gitops directory
func Add(o *AddParameters) error {
   if !o.SkipChecks {
       installed, err := checkTektonInstall()
       if err != nil {
           return fmt.Errorf("failed to run Tekton Pipelines installation check: %w", err)
       }
       if !installed {
           return errors.New("failed due to Tekton Pipelines or Triggers are not installed")
       }
   }
 
   gitopsName := getGitopsRepoName(o.GitopsRepo)
   gitopsPath := filepath.Join(o.Output, gitopsName)
   // log.Println(gitopsName, "+", gitopsPath, "=")
 
   exists, _ := isExisting(gitopsPath)
   if !exists {
       return fmt.Errorf("%s not exists at %s", gitopsName, gitopsPath)
   }
 
   // outputs := map[string]interface{}{}
 
   // sort.Strings(fileNames)
   // kustomize file should refer all the pipeline resources
   if err := addKustomize("bases", []string{"../base"}, filepath.Join(gitopsPath, appsDir, o.AppGitRepo, overlays, kustomize)); err != nil {
       return err
   }
 
   if err := addKustomize("bases", []string{"- ../../../services/service-1/overlays"}, filepath.Join(gitopsPath, appsDir, o.AppGitRepo, baseDir, kustomize)); err != nil {
       return err
   }
 
   return nil
}
 
func getAppDir(path, prefix string) string {
   return filepath.Join(path, envsDir, addPrefix(prefix, appsDir))
}
 

