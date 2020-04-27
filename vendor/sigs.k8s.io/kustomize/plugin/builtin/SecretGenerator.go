// Code generated by pluginator on SecretGenerator; DO NOT EDIT.
package builtin

import (
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/types"
	"sigs.k8s.io/yaml"
)

type SecretGeneratorPlugin struct {
	ldr              ifc.Loader
	rf               *resmap.Factory
	types.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	types.GeneratorOptions
	types.SecretArgs
}

func (p *SecretGeneratorPlugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, config []byte) (err error) {
	p.GeneratorOptions = types.GeneratorOptions{}
	p.SecretArgs = types.SecretArgs{}
	err = yaml.Unmarshal(config, p)
	if p.SecretArgs.Name == "" {
		p.SecretArgs.Name = p.Name
	}
	if p.SecretArgs.Namespace == "" {
		p.SecretArgs.Namespace = p.Namespace
	}
	p.ldr = ldr
	p.rf = rf
	return
}

func (p *SecretGeneratorPlugin) Generate() (resmap.ResMap, error) {
	return p.rf.FromSecretArgs(p.ldr, &p.GeneratorOptions, p.SecretArgs)
}

func NewSecretGeneratorPlugin() resmap.GeneratorPlugin {
	return &SecretGeneratorPlugin{}
}