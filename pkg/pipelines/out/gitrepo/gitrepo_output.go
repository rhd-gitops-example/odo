package gitrepo

import (
	"context"
	"os"
	"time"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/openshift/odo/pkg/pipelines/out"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/util"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type output struct {
	out.BaseOutput
	url    string
	user   string
	repo   string
	branch string
	token  string
}

// New creates an Output that outputs to a Git repository
func New(url, user, repo, branch, token string) (out.Output, error) {
	return &output{
		BaseOutput: out.New(),
		branch:     branch,
		url:        url,
		token:      token,
		repo:       repo,
		user:       user,
	}, nil
}

// Write all outpout items to a Git repository
func (o *output) Write() error {
	err := o.createRepo()
	if err != nil {
		return err
	}

	err = o.initRepo()
	if err != nil {
		return err
	}

	err = o.createBranch()
	if err != nil {
		return err
	}

	o.createPullRequest()
	if err != nil {
		return err
	}

	return nil
}

func (o *output) createRepo() error {
	client, err := factory.NewClient("github", "", o.token)
	if err != nil {
		return err
	}

	in := &scm.RepositoryInput{
		Name:        o.repo,
		Description: "GitOps Test Repository",
	}

	_, _, err = client.Repositories.Create(context.Background(), in)
	return err
}

func (o *output) initRepo() error {
	fs := memfs.New()
	r, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	util.WriteFile(fs, "README.md", []byte("This is my test repo"), 0644)

	_, err = w.Add("README.md")
	if err != nil {
		return err
	}

	_, err = w.Commit("This is a test commit comment\n", &git.CommitOptions{Author: defaultSignature()})
	if err != nil {
		return err
	}

	remote, err := r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{o.url},
	})
	if err != nil {
		return err
	}

	err = remote.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "test",
			Password: o.token,
		},
	})

	return err
}

func (o *output) createBranch() error {
	fs := memfs.New()
	r, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: o.url,
		Auth: &http.BasicAuth{
			Password: o.token,
		},
	})
	if err != nil {
		return err
	}

	headRef, err := r.Head()
	if err != nil {
		return err
	}

	ref := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+o.branch), headRef.Hash())
	err = r.Storer.SetReference(ref)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = o.populateWorkTree(w, fs)
	if err != nil {
		return err
	}

	_, err = w.Commit("This is a test commit comment\n", &git.CommitOptions{Author: defaultSignature()})
	if err != nil {
		return err
	}

	return r.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("refs/heads/master:refs/heads/" + o.branch),
			":refs/heads/branch",
		},
		Auth: &http.BasicAuth{
			Username: "test",
			Password: o.token,
		},
		Prune: true,
	})
}

func (o *output) populateWorkTree(w *git.Worktree, fs billy.Filesystem) error {
	for path, data := range o.Items {
		err := o.write(fs, path, data)
		if err != nil {
			return err
		}
		_, err = w.Add(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *output) write(fs billy.Filesystem, path string, data interface{}) error {
	f, err := fs.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = out.Marshal(f, data)
	if err != nil {
		return err
	}
	return nil
}

func (o *output) createPullRequest() error {
	client, err := factory.NewClient("github", "", o.token)
	if err != nil {
		return err
	}

	input := &scm.PullRequestInput{
		Title: "GitOps Test PR",
		Body:  "Please pull these awesome changes",
		Head:  o.user + ":" + o.branch,
		Base:  "master",
	}

	_, _, err = client.PullRequests.Create(context.Background(), o.user+"/"+o.repo, input)
	return err
}

func defaultSignature() *object.Signature {
	when, _ := time.Parse(object.DateFormat, "Thu May 04 00:03:43 2019 +0200")
	return &object.Signature{
		Name:  "Test User",
		Email: "testuser@redhat.com",
		When:  when,
	}
}
