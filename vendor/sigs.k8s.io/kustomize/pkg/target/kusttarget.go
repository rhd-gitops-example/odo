<<<<<<< HEAD
/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package target implements state for the set of all resources to customize.
=======
// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// Package target implements state for the set of all
// resources to customize.
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
package target

import (
	"bytes"
	"encoding/json"
	"fmt"
<<<<<<< HEAD
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/pkg/constants"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/ifc/transformer"
	interror "sigs.k8s.io/kustomize/pkg/internal/error"
	patchtransformer "sigs.k8s.io/kustomize/pkg/patch/transformer"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/transformers"
	"sigs.k8s.io/kustomize/pkg/transformers/config"
	"sigs.k8s.io/kustomize/pkg/types"
=======
	"log"
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/v3/pkg/accumulator"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/pgmconfig"
	"sigs.k8s.io/kustomize/v3/pkg/plugins"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/transformers"
	"sigs.k8s.io/kustomize/v3/pkg/transformers/config"
	"sigs.k8s.io/kustomize/v3/pkg/types"
	"sigs.k8s.io/kustomize/v3/plugin/builtin"
	"sigs.k8s.io/yaml"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

// KustTarget encapsulates the entirety of a kustomization build.
type KustTarget struct {
	kustomization *types.Kustomization
	ldr           ifc.Loader
	rFactory      *resmap.Factory
<<<<<<< HEAD
	tFactory      transformer.Factory
=======
	tFactory      resmap.PatchFactory
	pLdr          *plugins.Loader
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// NewKustTarget returns a new instance of KustTarget primed with a Loader.
func NewKustTarget(
	ldr ifc.Loader,
	rFactory *resmap.Factory,
<<<<<<< HEAD
	tFactory transformer.Factory) (*KustTarget, error) {
=======
	tFactory resmap.PatchFactory,
	pLdr *plugins.Loader) (*KustTarget, error) {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	content, err := loadKustFile(ldr)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD
	content = types.DealWithDeprecatedFields(content)
=======
	content = types.FixKustomizationPreUnmarshalling(content)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	var k types.Kustomization
	err = unmarshal(content, &k)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD
	errs := k.EnforceFields()
	if len(errs) > 0 {
		return nil, fmt.Errorf("Failed to read kustomization file under %s:\n"+strings.Join(errs, "\n"), ldr.Root())
=======
	k.FixKustomizationPostUnmarshalling()
	errs := k.EnforceFields()
	if len(errs) > 0 {
		return nil, fmt.Errorf(
			"Failed to read kustomization file under %s:\n"+
				strings.Join(errs, "\n"), ldr.Root())
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
	return &KustTarget{
		kustomization: &k,
		ldr:           ldr,
		rFactory:      rFactory,
		tFactory:      tFactory,
<<<<<<< HEAD
=======
		pLdr:          pLdr,
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}, nil
}

func quoted(l []string) []string {
	r := make([]string, len(l))
	for i, v := range l {
		r[i] = "'" + v + "'"
	}
	return r
}

func commaOr(q []string) string {
	return strings.Join(q[:len(q)-1], ", ") + " or " + q[len(q)-1]
}

func loadKustFile(ldr ifc.Loader) ([]byte, error) {
	var content []byte
	match := 0
<<<<<<< HEAD
	for _, kf := range constants.KustomizationFileNames {
=======
	for _, kf := range pgmconfig.RecognizedKustomizationFileNames() {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		c, err := ldr.Load(kf)
		if err == nil {
			match += 1
			content = c
		}
	}
	switch match {
	case 0:
		return nil, fmt.Errorf(
			"unable to find one of %v in directory '%s'",
<<<<<<< HEAD
			commaOr(quoted(constants.KustomizationFileNames)), ldr.Root())
	case 1:
		return content, nil
	default:
		return nil, fmt.Errorf("Found multiple kustomization files under: %s\n", ldr.Root())
=======
			commaOr(quoted(pgmconfig.RecognizedKustomizationFileNames())),
			ldr.Root())
	case 1:
		return content, nil
	default:
		return nil, fmt.Errorf(
			"Found multiple kustomization files under: %s\n", ldr.Root())
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
}

func unmarshal(y []byte, o interface{}) error {
	j, err := yaml.YAMLToJSON(y)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(j))
	dec.DisallowUnknownFields()
	return dec.Decode(o)
}

// MakeCustomizedResMap creates a ResMap per kustomization instructions.
// The Resources in the returned ResMap are fully customized.
func (kt *KustTarget) MakeCustomizedResMap() (resmap.ResMap, error) {
<<<<<<< HEAD
=======
	return kt.makeCustomizedResMap(types.GarbageIgnore)
}

func (kt *KustTarget) MakePruneConfigMap() (resmap.ResMap, error) {
	return kt.makeCustomizedResMap(types.GarbageCollect)
}

func (kt *KustTarget) makeCustomizedResMap(
	garbagePolicy types.GarbagePolicy) (resmap.ResMap, error) {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	ra, err := kt.AccumulateTarget()
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD
	err = ra.Transform(kt.tFactory.MakeHashTransformer())
	if err != nil {
		return nil, err
	}
=======

	// The following steps must be done last, not as part of
	// the recursion implicit in AccumulateTarget.

	err = kt.addHashesToNames(ra)
	if err != nil {
		return nil, err
	}

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	// Given that names have changed (prefixs/suffixes added),
	// fix all the back references to those names.
	err = ra.FixBackReferences()
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD
	// With all the back references fixed, it's OK to resolve Vars.
	err = ra.ResolveVars()
	return ra.ResMap(), err
=======

	// With all the back references fixed, it's OK to resolve Vars.
	err = ra.ResolveVars()
	if err != nil {
		return nil, err
	}

	err = kt.computeInventory(ra, garbagePolicy)
	if err != nil {
		return nil, err
	}

	return ra.ResMap(), nil
}

func (kt *KustTarget) addHashesToNames(
	ra *accumulator.ResAccumulator) error {
	p := builtin.NewHashTransformerPlugin()
	err := kt.configureBuiltinPlugin(p, nil, plugins.HashTransformer)
	if err != nil {
		return err
	}
	return ra.Transform(p)
}

func (kt *KustTarget) computeInventory(
	ra *accumulator.ResAccumulator, garbagePolicy types.GarbagePolicy) error {
	inv := kt.kustomization.Inventory
	if inv == nil {
		return nil
	}
	if inv.Type != "ConfigMap" {
		return fmt.Errorf("don't know how to do that")
	}

	if inv.ConfigMap.Namespace != kt.kustomization.Namespace {
		return fmt.Errorf("namespace mismatch")
	}

	var c struct {
		Policy           string
		types.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	}
	c.Name = inv.ConfigMap.Name
	c.Namespace = inv.ConfigMap.Namespace
	c.Policy = garbagePolicy.String()
	p := builtin.NewInventoryTransformerPlugin()
	err := kt.configureBuiltinPlugin(p, c, plugins.InventoryTransformer)
	if err != nil {
		return err
	}
	return ra.Transform(p)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

func (kt *KustTarget) shouldAddHashSuffixesToGeneratedResources() bool {
	return kt.kustomization.GeneratorOptions == nil ||
		!kt.kustomization.GeneratorOptions.DisableNameSuffixHash
}

// AccumulateTarget returns a new ResAccumulator,
// holding customized resources and the data/rules used
// to do so.  The name back references and vars are
// not yet fixed.
func (kt *KustTarget) AccumulateTarget() (
<<<<<<< HEAD
	ra *ResAccumulator, err error) {
	// TODO(monopole): Get rid of the KustomizationErrors accumulator.
	// It's not consistently used, and complicates tests.
	errs := &interror.KustomizationErrors{}
	ra, errs = kt.accumulateBases()
	resources, err := kt.rFactory.FromFiles(
		kt.ldr, kt.kustomization.Resources)
	if err != nil {
		errs.Append(errors.Wrap(err, "rawResources failed to read Resources"))
	}
	if len(errs.Get()) > 0 {
		return ra, errs
	}
	err = ra.MergeResourcesWithErrorOnIdCollision(resources)
	if err != nil {
		errs.Append(errors.Wrap(err, "MergeResourcesWithErrorOnIdCollision"))
=======
	ra *accumulator.ResAccumulator, err error) {
	ra = accumulator.MakeEmptyAccumulator()
	err = kt.accumulateResources(ra, kt.kustomization.Resources)
	if err != nil {
		return nil, errors.Wrap(err, "accumulating resources")
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
	tConfig, err := config.MakeTransformerConfig(
		kt.ldr, kt.kustomization.Configurations)
	if err != nil {
		return nil, err
	}
	err = ra.MergeConfig(tConfig)
	if err != nil {
<<<<<<< HEAD
		errs.Append(errors.Wrap(err, "MergeConfig"))
	}
	err = ra.MergeVars(kt.kustomization.Vars)
	if err != nil {
		errs.Append(errors.Wrap(err, "MergeVars"))
	}
	crdTc, err := config.LoadConfigFromCRDs(kt.ldr, kt.kustomization.Crds)
	if err != nil {
		errs.Append(errors.Wrap(err, "LoadCRDs"))
	}
	err = ra.MergeConfig(crdTc)
	if err != nil {
		errs.Append(errors.Wrap(err, "merge CRDs"))
	}
	resMap, err := kt.generateConfigMapsAndSecrets(errs)
	if err != nil {
		errs.Append(errors.Wrap(err, "generateConfigMapsAndSecrets"))
	}
	err = ra.MergeResourcesWithOverride(resMap)
	if err != nil {
		return nil, err
	}
	patches, err := kt.rFactory.RF().SliceFromPatches(
		kt.ldr, kt.kustomization.PatchesStrategicMerge)
	if err != nil {
		errs.Append(errors.Wrap(err, "SliceFromPatches"))
	}
	if len(errs.Get()) > 0 {
		return nil, errs
	}
	t, err := kt.newTransformer(patches, ra.tConfig)
	if err != nil {
		return nil, err
	}
	err = ra.Transform(t)
	if err != nil {
		return nil, err
	}
	return ra, nil
}

func (kt *KustTarget) generateConfigMapsAndSecrets(
	errs *interror.KustomizationErrors) (resmap.ResMap, error) {
	kt.rFactory.Set(kt.ldr)
	cms, err := kt.rFactory.NewResMapFromConfigMapArgs(
		kt.kustomization.ConfigMapGenerator, kt.kustomization.GeneratorOptions)
	if err != nil {
		errs.Append(errors.Wrap(err, "NewResMapFromConfigMapArgs"))
	}
	secrets, err := kt.rFactory.NewResMapFromSecretArgs(
		kt.kustomization.SecretGenerator, kt.kustomization.GeneratorOptions)
	if err != nil {
		errs.Append(errors.Wrap(err, "NewResMapFromSecretArgs"))
	}
	return resmap.MergeWithErrorOnIdCollision(cms, secrets)
}

// accumulateBases returns a new ResAccumulator
// holding customized resources and the data/rules
// used to customized them from only the _bases_
// of this KustTarget.
func (kt *KustTarget) accumulateBases() (
	ra *ResAccumulator, errs *interror.KustomizationErrors) {
	errs = &interror.KustomizationErrors{}
	ra = MakeEmptyAccumulator()

	for _, path := range kt.kustomization.Bases {
		ldr, err := kt.ldr.New(path)
		if err != nil {
			errs.Append(errors.Wrap(err, "couldn't make loader for "+path))
			continue
		}
		subKt, err := NewKustTarget(
			ldr, kt.rFactory, kt.tFactory)
		if err != nil {
			errs.Append(errors.Wrap(err, "couldn't make target for "+path))
			ldr.Cleanup()
			continue
		}
		subRa, err := subKt.AccumulateTarget()
		if err != nil {
			errs.Append(errors.Wrap(err, "AccumulateTarget"))
			ldr.Cleanup()
			continue
		}
		err = ra.MergeAccumulator(subRa)
		if err != nil {
			errs.Append(errors.Wrap(err, path))
		}
		ldr.Cleanup()
	}
	return ra, errs
}

// newTransformer makes a Transformer that does a collection
// of object transformations.
func (kt *KustTarget) newTransformer(
	patches []*resource.Resource, tConfig *config.TransformerConfig) (
	transformers.Transformer, error) {
	var r []transformers.Transformer
	t, err := kt.tFactory.MakePatchTransformer(patches, kt.rFactory.RF())
	if err != nil {
		return nil, err
	}
	r = append(r, t)
	r = append(r, transformers.NewNamespaceTransformer(
		string(kt.kustomization.Namespace), tConfig.NameSpace))
	t, err = transformers.NewNamePrefixSuffixTransformer(
		string(kt.kustomization.NamePrefix),
		string(kt.kustomization.NameSuffix),
		tConfig.NamePrefix,
	)
	if err != nil {
		return nil, err
	}
	r = append(r, t)
	t, err = transformers.NewLabelsMapTransformer(
		kt.kustomization.CommonLabels, tConfig.CommonLabels)
	if err != nil {
		return nil, err
	}
	r = append(r, t)
	t, err = transformers.NewAnnotationsMapTransformer(
		kt.kustomization.CommonAnnotations, tConfig.CommonAnnotations)
	if err != nil {
		return nil, err
	}
	r = append(r, t)
	t, err = patchtransformer.NewPatchJson6902Factory(kt.ldr).
		MakePatchJson6902Transformer(kt.kustomization.PatchesJson6902)
	if err != nil {
		return nil, err
	}
	r = append(r, t)
	t, err = transformers.NewImageTransformer(kt.kustomization.Images)
	if err != nil {
		return nil, err
	}
	r = append(r, t)
	return transformers.NewMultiTransformer(r), nil
=======
		return nil, errors.Wrapf(
			err, "merging config %v", tConfig)
	}
	crdTc, err := config.LoadConfigFromCRDs(kt.ldr, kt.kustomization.Crds)
	if err != nil {
		return nil, errors.Wrapf(
			err, "loading CRDs %v", kt.kustomization.Crds)
	}
	err = ra.MergeConfig(crdTc)
	if err != nil {
		return nil, errors.Wrapf(
			err, "merging CRDs %v", crdTc)
	}
	err = kt.runGenerators(ra)
	if err != nil {
		return nil, err
	}
	err = kt.runTransformers(ra)
	if err != nil {
		return nil, err
	}
	err = ra.MergeVars(kt.kustomization.Vars)
	if err != nil {
		return nil, errors.Wrapf(
			err, "merging vars %v", kt.kustomization.Vars)
	}
	return ra, nil
}

func (kt *KustTarget) runGenerators(
	ra *accumulator.ResAccumulator) error {
	var generators []resmap.Generator
	gs, err := kt.configureBuiltinGenerators()
	if err != nil {
		return err
	}
	generators = append(generators, gs...)
	gs, err = kt.configureExternalGenerators()
	if err != nil {
		return errors.Wrap(err, "loading generator plugins")
	}
	generators = append(generators, gs...)
	for _, g := range generators {
		resMap, err := g.Generate()
		if err != nil {
			return err
		}
		err = ra.AbsorbAll(resMap)
		if err != nil {
			return errors.Wrapf(err, "merging from generator %v", g)
		}
	}
	return nil
}

func (kt *KustTarget) configureExternalGenerators() ([]resmap.Generator, error) {
	ra := accumulator.MakeEmptyAccumulator()
	err := kt.accumulateResources(ra, kt.kustomization.Generators)
	if err != nil {
		return nil, err
	}
	return kt.pLdr.LoadGenerators(kt.ldr, ra.ResMap())
}

func (kt *KustTarget) runTransformers(ra *accumulator.ResAccumulator) error {
	var r []resmap.Transformer
	tConfig := ra.GetTransformerConfig()
	lts, err := kt.configureBuiltinTransformers(tConfig)
	if err != nil {
		return err
	}
	r = append(r, lts...)
	lts, err = kt.configureExternalTransformers()
	if err != nil {
		return err
	}
	r = append(r, lts...)
	t := transformers.NewMultiTransformer(r)
	return ra.Transform(t)
}

func (kt *KustTarget) configureExternalTransformers() ([]resmap.Transformer, error) {
	ra := accumulator.MakeEmptyAccumulator()
	err := kt.accumulateResources(ra, kt.kustomization.Transformers)
	if err != nil {
		return nil, err
	}
	return kt.pLdr.LoadTransformers(kt.ldr, ra.ResMap())
}

// accumulateResources fills the given resourceAccumulator
// with resources read from the given list of paths.
func (kt *KustTarget) accumulateResources(
	ra *accumulator.ResAccumulator, paths []string) error {
	for _, path := range paths {
		ldr, err := kt.ldr.New(path)
		if err == nil {
			err = kt.accumulateDirectory(ra, ldr, path)
			if err != nil {
				return err
			}
		} else {
			err2 := kt.accumulateFile(ra, path)
			if err2 != nil {
				// Log ldr.New() error to highlight git failures.
				log.Print(err.Error())
				return err2
			}
		}
	}
	return nil
}

func (kt *KustTarget) accumulateDirectory(
	ra *accumulator.ResAccumulator, ldr ifc.Loader, path string) error {
	defer ldr.Cleanup()
	subKt, err := NewKustTarget(
		ldr, kt.rFactory, kt.tFactory, kt.pLdr)
	if err != nil {
		return errors.Wrapf(err, "couldn't make target for path '%s'", path)
	}
	subRa, err := subKt.AccumulateTarget()
	if err != nil {
		return errors.Wrapf(
			err, "recursed accumulation of path '%s'", path)
	}
	err = ra.MergeAccumulator(subRa)
	if err != nil {
		return errors.Wrapf(
			err, "recursed merging from path '%s'", path)
	}
	return nil
}

func (kt *KustTarget) accumulateFile(
	ra *accumulator.ResAccumulator, path string) error {
	resources, err := kt.rFactory.FromFile(kt.ldr, path)
	if err != nil {
		return errors.Wrapf(err, "accumulating resources from '%s'", path)
	}
	err = ra.AppendAll(resources)
	if err != nil {
		return errors.Wrapf(err, "merging resources from '%s'", path)
	}
	return nil
}

func (kt *KustTarget) configureBuiltinPlugin(
	p resmap.Configurable, c interface{}, bpt plugins.BuiltinPluginType) (err error) {
	var y []byte
	if c != nil {
		y, err = yaml.Marshal(c)
		if err != nil {
			return errors.Wrapf(
				err, "builtin %s marshal", bpt)
		}
	}
	err = p.Config(kt.ldr, kt.rFactory, y)
	if err != nil {
		return errors.Wrapf(err, "builtin %s config: %v", bpt, y)
	}
	return nil
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}
