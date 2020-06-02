package pipelines

import (
	"bytes"
	"fmt"
	"text/template"
)

const scriptTemplate = `#!/bin/sh
is_argocd=false
argo_path="config/argocd"
cicd_path="config/{{ .CICDEnv }}"
client={{ .Client }}
if [ -d ${argo_path} ]
then
   printf "Apply $(basename ${argo_path}) applications\n"
   if [ "$client" != "" ];then $client apply --dry-run=$(inputs.params.DRYRUN) -k "${argo_path}/config"; fi
   is_argocd=true
fi
printf "Apply $(basename ${cicd_path}) environment\n"
if [ "$client" != "" ];then $client apply --dry-run=$(inputs.params.DRYRUN) -k "${cicd_path}/overlays"; fi
for dir in $(ls -d environments/*/)
do
   	if ! $is_argocd
   	then
	   printf "Apply $(basename ${dir}) environment\n"
       env_path="${dir}env/overlays"
       if [ "$client" != "" ];then $client apply --dry-run=$(inputs.params.DRYRUN) -k $env_path; fi
   	else
       if [ -d "${dir}apps" ]
       then
           for app in $(ls -d ${dir}apps/*/)
           do
            	printf "Apply $(basename ${app}) application\n"
            	if [ "$client" != "" ];then $client apply --dry-run=$(inputs.params.DRYRUN) -k $app; fi
           done
       else 
        	printf "Apply $(basename ${dir}) environment\n"
            env_path="${dir}env/overlays"
            if [ "$client" != "" ];then $client apply --dry-run=$(inputs.params.DRYRUN) -k $env_path; fi
       fi
   	fi
done`

type templateParam struct {
	Client  string
	CICDEnv string
}

func makeScript(client, cicdEnv string) (string, error) {
	params := templateParam{CICDEnv: cicdEnv, Client: client}
	template, err := template.New("dryrun_script").Parse(scriptTemplate)
	if err != nil {
		return "", fmt.Errorf("unable to parse template: %v", err)
	}
	var buf bytes.Buffer
	err = template.Execute(&buf, params)
	if err != nil {
		return "", fmt.Errorf("unable to execute template: %v", err)
	}
	return buf.String(), nil
}
