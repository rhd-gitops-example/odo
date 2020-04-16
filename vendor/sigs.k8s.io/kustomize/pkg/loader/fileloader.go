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
=======
// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)

package loader

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

<<<<<<< HEAD
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/git"
	"sigs.k8s.io/kustomize/pkg/ifc"
=======
	"sigs.k8s.io/kustomize/v3/pkg/fs"
	"sigs.k8s.io/kustomize/v3/pkg/git"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

// fileLoader is a kustomization's interface to files.
//
// The directory in which a kustomization file sits
<<<<<<< HEAD
// is referred to below as the kustomization's root.
=======
// is referred to below as the kustomization's _root_.
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
//
// An instance of fileLoader has an immutable root,
// and offers a `New` method returning a new loader
// with a new root.
//
// A kustomization file refers to two kinds of files:
//
// * supplemental data paths
//
//   `Load` is used to visit these paths.
//
<<<<<<< HEAD
//   They must terminate in or below the root.
//
//   They hold things like resources, patches,
//   data for ConfigMaps, etc.
//
// * bases; other kustomizations
=======
//   These paths refer to resources, patches,
//   data for ConfigMaps and Secrets, etc.
//
//   The loadRestrictor may disallow certain paths
//   or classes of paths.
//
// * bases (other kustomizations)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
//
//   `New` is used to load bases.
//
//   A base can be either a remote git repo URL, or
//   a directory specified relative to the current
//   root. In the former case, the repo is locally
//   cloned, and the new loader is rooted on a path
//   in that clone.
//
//   As loaders create new loaders, a root history
//   is established, and used to disallow:
//
//   - A base that is a repository that, in turn,
//     specifies a base repository seen previously
//     in the loading stack (a cycle).
//
//   - An overlay depending on a base positioned at
//     or above it.  I.e. '../foo' is OK, but '.',
//     '..', '../..', etc. are disallowed.  Allowing
//     such a base has no advantages and encourages
//     cycles, particularly if some future change
//     were to introduce globbing to file
//     specifications in the kustomization file.
//
// These restrictions assure that kustomizations
// are self-contained and relocatable, and impose
// some safety when relying on remote kustomizations,
<<<<<<< HEAD
// e.g. a ConfigMap generator specified to read
// from /etc/passwd will fail.
=======
// e.g. a remotely loaded ConfigMap generator specified
// to read from /etc/passwd will fail.
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
//
type fileLoader struct {
	// Loader that spawned this loader.
	// Used to avoid cycles.
	referrer *fileLoader
<<<<<<< HEAD
	// An absolute, cleaned path to a directory.
	// The Load function reads from this directory,
	// or directories below it.
	root fs.ConfirmedDir
	// If this is non-nil, the files were
	// obtained from the given repository.
	repoSpec *git.RepoSpec
	// File system utilities.
	fSys fs.FileSystem
	// Used to clone repositories.
	cloner git.Cloner
=======

	// An absolute, cleaned path to a directory.
	// The Load function will read non-absolute
	// paths relative to this directory.
	root fs.ConfirmedDir

	// Restricts behavior of Load function.
	loadRestrictor LoadRestrictorFunc

	// Used to validate various k8s data fields.
	validator ifc.Validator

	// If this is non-nil, the files were
	// obtained from the given repository.
	repoSpec *git.RepoSpec

	// File system utilities.
	fSys fs.FileSystem

	// Used to clone repositories.
	cloner git.Cloner

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	// Used to clean up, as needed.
	cleaner func() error
}

<<<<<<< HEAD
// NewFileLoaderAtCwd returns a loader that loads from ".".
func NewFileLoaderAtCwd(fSys fs.FileSystem) *fileLoader {
	return newLoaderOrDie(fSys, ".")
}

// NewFileLoaderAtRoot returns a loader that loads from "/".
func NewFileLoaderAtRoot(fSys fs.FileSystem) *fileLoader {
	return newLoaderOrDie(fSys, string(filepath.Separator))
=======
const CWD = "."

// NewFileLoaderAtCwd returns a loader that loads from ".".
// A convenience for kustomize edit commands.
func NewFileLoaderAtCwd(v ifc.Validator, fSys fs.FileSystem) *fileLoader {
	return newLoaderOrDie(
		RestrictionRootOnly, v, fSys, CWD)
}

// NewFileLoaderAtRoot returns a loader that loads from "/".
// A convenience for tests.
func NewFileLoaderAtRoot(v ifc.Validator, fSys fs.FileSystem) *fileLoader {
	return newLoaderOrDie(
		RestrictionRootOnly, v, fSys, string(filepath.Separator))
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// Root returns the absolute path that is prepended to any
// relative paths used in Load.
<<<<<<< HEAD
func (l *fileLoader) Root() string {
	return l.root.String()
}

func newLoaderOrDie(fSys fs.FileSystem, path string) *fileLoader {
=======
func (fl *fileLoader) Root() string {
	return fl.root.String()
}

func newLoaderOrDie(
	lr LoadRestrictorFunc, v ifc.Validator,
	fSys fs.FileSystem, path string) *fileLoader {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	root, err := demandDirectoryRoot(fSys, path)
	if err != nil {
		log.Fatalf("unable to make loader at '%s'; %v", path, err)
	}
	return newLoaderAtConfirmedDir(
<<<<<<< HEAD
		root, fSys, nil, git.ClonerUsingGitExec)
=======
		lr, v, root, fSys, nil, git.ClonerUsingGitExec)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// newLoaderAtConfirmedDir returns a new fileLoader with given root.
func newLoaderAtConfirmedDir(
<<<<<<< HEAD
	root fs.ConfirmedDir, fSys fs.FileSystem,
	referrer *fileLoader, cloner git.Cloner) *fileLoader {
	return &fileLoader{
		root:     root,
		referrer: referrer,
		fSys:     fSys,
		cloner:   cloner,
		cleaner:  func() error { return nil },
=======
	lr LoadRestrictorFunc,
	v ifc.Validator,
	root fs.ConfirmedDir, fSys fs.FileSystem,
	referrer *fileLoader, cloner git.Cloner) *fileLoader {
	return &fileLoader{
		loadRestrictor: lr,
		validator:      v,
		root:           root,
		referrer:       referrer,
		fSys:           fSys,
		cloner:         cloner,
		cleaner:        func() error { return nil },
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
}

// Assure that the given path is in fact a directory.
func demandDirectoryRoot(
	fSys fs.FileSystem, path string) (fs.ConfirmedDir, error) {
	if path == "" {
		return "", fmt.Errorf(
			"loader root cannot be empty")
	}
	d, f, err := fSys.CleanedAbs(path)
	if err != nil {
		return "", fmt.Errorf(
			"absolute path error in '%s' : %v", path, err)
	}
	if f != "" {
		return "", fmt.Errorf(
			"got file '%s', but '%s' must be a directory to be a root",
			f, path)
	}
	return d, nil
}

// New returns a new Loader, rooted relative to current loader,
// or rooted in a temp directory holding a git repo clone.
<<<<<<< HEAD
func (l *fileLoader) New(path string) (ifc.Loader, error) {
=======
func (fl *fileLoader) New(path string) (ifc.Loader, error) {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	if path == "" {
		return nil, fmt.Errorf("new root cannot be empty")
	}
	repoSpec, err := git.NewRepoSpecFromUrl(path)
	if err == nil {
		// Treat this as git repo clone request.
<<<<<<< HEAD
		if err := l.errIfRepoCycle(repoSpec); err != nil {
			return nil, err
		}
		return newLoaderAtGitClone(repoSpec, l.fSys, l.referrer, l.cloner)
=======
		if err := fl.errIfRepoCycle(repoSpec); err != nil {
			return nil, err
		}
		return newLoaderAtGitClone(
			repoSpec, fl.validator, fl.fSys, fl, fl.cloner)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
	if filepath.IsAbs(path) {
		return nil, fmt.Errorf("new root '%s' cannot be absolute", path)
	}
<<<<<<< HEAD
	root, err := demandDirectoryRoot(l.fSys, l.root.Join(path))
	if err != nil {
		return nil, err
	}
	if err := l.errIfGitContainmentViolation(root); err != nil {
		return nil, err
	}
	if err := l.errIfArgEqualOrHigher(root); err != nil {
		return nil, err
	}
	return newLoaderAtConfirmedDir(
		root, l.fSys, l, l.cloner), nil
=======
	root, err := demandDirectoryRoot(fl.fSys, fl.root.Join(path))
	if err != nil {
		return nil, err
	}
	if err := fl.errIfGitContainmentViolation(root); err != nil {
		return nil, err
	}
	if err := fl.errIfArgEqualOrHigher(root); err != nil {
		return nil, err
	}
	return newLoaderAtConfirmedDir(
		fl.loadRestrictor, fl.validator, root, fl.fSys, fl, fl.cloner), nil
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// newLoaderAtGitClone returns a new Loader pinned to a temporary
// directory holding a cloned git repo.
func newLoaderAtGitClone(
<<<<<<< HEAD
	repoSpec *git.RepoSpec, fSys fs.FileSystem,
=======
	repoSpec *git.RepoSpec,
	v ifc.Validator, fSys fs.FileSystem,
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	referrer *fileLoader, cloner git.Cloner) (ifc.Loader, error) {
	err := cloner(repoSpec)
	if err != nil {
		return nil, err
	}
	root, f, err := fSys.CleanedAbs(repoSpec.AbsPath())
	if err != nil {
		return nil, err
	}
	// We don't know that the path requested in repoSpec
	// is a directory until we actually clone it and look
	// inside.  That just happened, hence the error check
	// is here.
	if f != "" {
		return nil, fmt.Errorf(
			"'%s' refers to file '%s'; expecting directory",
			repoSpec.AbsPath(), f)
	}
	return &fileLoader{
<<<<<<< HEAD
		root:     root,
		referrer: referrer,
		repoSpec: repoSpec,
		fSys:     fSys,
		cloner:   cloner,
		cleaner:  repoSpec.Cleaner(fSys),
	}, nil
}

func (l *fileLoader) errIfGitContainmentViolation(
	base fs.ConfirmedDir) error {
	containingRepo := l.containingRepo()
=======
		// Clones never allowed to escape root.
		loadRestrictor: RestrictionRootOnly,
		validator:      v,
		root:           root,
		referrer:       referrer,
		repoSpec:       repoSpec,
		fSys:           fSys,
		cloner:         cloner,
		cleaner:        repoSpec.Cleaner(fSys),
	}, nil
}

func (fl *fileLoader) errIfGitContainmentViolation(
	base fs.ConfirmedDir) error {
	containingRepo := fl.containingRepo()
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	if containingRepo == nil {
		return nil
	}
	if !base.HasPrefix(containingRepo.CloneDir()) {
		return fmt.Errorf(
			"security; bases in kustomizations found in "+
				"cloned git repos must be within the repo, "+
				"but base '%s' is outside '%s'",
			base, containingRepo.CloneDir())
	}
	return nil
}

// Looks back through referrers for a git repo, returning nil
// if none found.
<<<<<<< HEAD
func (l *fileLoader) containingRepo() *git.RepoSpec {
	if l.repoSpec != nil {
		return l.repoSpec
	}
	if l.referrer == nil {
		return nil
	}
	return l.referrer.containingRepo()
=======
func (fl *fileLoader) containingRepo() *git.RepoSpec {
	if fl.repoSpec != nil {
		return fl.repoSpec
	}
	if fl.referrer == nil {
		return nil
	}
	return fl.referrer.containingRepo()
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// errIfArgEqualOrHigher tests whether the argument,
// is equal to or above the root of any ancestor.
<<<<<<< HEAD
func (l *fileLoader) errIfArgEqualOrHigher(
	candidateRoot fs.ConfirmedDir) error {
	if l.root.HasPrefix(candidateRoot) {
		return fmt.Errorf(
			"cycle detected: candidate root '%s' contains visited root '%s'",
			candidateRoot, l.root)
	}
	if l.referrer == nil {
		return nil
	}
	return l.referrer.errIfArgEqualOrHigher(candidateRoot)
=======
func (fl *fileLoader) errIfArgEqualOrHigher(
	candidateRoot fs.ConfirmedDir) error {
	if fl.root.HasPrefix(candidateRoot) {
		return fmt.Errorf(
			"cycle detected: candidate root '%s' contains visited root '%s'",
			candidateRoot, fl.root)
	}
	if fl.referrer == nil {
		return nil
	}
	return fl.referrer.errIfArgEqualOrHigher(candidateRoot)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// TODO(monopole): Distinguish branches?
// I.e. Allow a distinction between git URI with
// path foo and tag bar and a git URI with the same
// path but a different tag?
<<<<<<< HEAD
func (l *fileLoader) errIfRepoCycle(newRepoSpec *git.RepoSpec) error {
	// TODO(monopole): Use parsed data instead of Raw().
	if l.repoSpec != nil &&
		strings.HasPrefix(l.repoSpec.Raw(), newRepoSpec.Raw()) {
		return fmt.Errorf(
			"cycle detected: URI '%s' referenced by previous URI '%s'",
			newRepoSpec.Raw(), l.repoSpec.Raw())
	}
	if l.referrer == nil {
		return nil
	}
	return l.referrer.errIfRepoCycle(newRepoSpec)
}

// Load returns content of file at the given relative path,
// else an error.  The path must refer to a file in or
// below the current root.
func (l *fileLoader) Load(path string) ([]byte, error) {
	if filepath.IsAbs(path) {
		return nil, l.loadOutOfBounds(path)
	}
	d, f, err := l.fSys.CleanedAbs(l.root.Join(path))
	if err != nil {
		return nil, err
	}
	if f == "" {
		return nil, fmt.Errorf(
			"'%s' must be a file (got d='%s')", path, d)
	}
	if !d.HasPrefix(l.root) {
		return nil, l.loadOutOfBounds(path)
	}
	return l.fSys.ReadFile(d.Join(f))
}

func (l *fileLoader) loadOutOfBounds(path string) error {
	return fmt.Errorf(
		"security; file '%s' is not in or below '%s'",
		path, l.root)
}

// Cleanup runs the cleaner.
func (l *fileLoader) Cleanup() error {
	return l.cleaner()
=======
func (fl *fileLoader) errIfRepoCycle(newRepoSpec *git.RepoSpec) error {
	// TODO(monopole): Use parsed data instead of Raw().
	if fl.repoSpec != nil &&
		strings.HasPrefix(fl.repoSpec.Raw(), newRepoSpec.Raw()) {
		return fmt.Errorf(
			"cycle detected: URI '%s' referenced by previous URI '%s'",
			newRepoSpec.Raw(), fl.repoSpec.Raw())
	}
	if fl.referrer == nil {
		return nil
	}
	return fl.referrer.errIfRepoCycle(newRepoSpec)
}

// Load returns the content of file at the given path,
// else an error.  Relative paths are taken relative
// to the root.
func (fl *fileLoader) Load(path string) ([]byte, error) {
	if !filepath.IsAbs(path) {
		path = fl.root.Join(path)
	}
	path, err := fl.loadRestrictor(fl.fSys, fl.root, path)
	if err != nil {
		return nil, err
	}
	return fl.fSys.ReadFile(path)
}

// Cleanup runs the cleaner.
func (fl *fileLoader) Cleanup() error {
	return fl.cleaner()
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}
