package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosimple/slug"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func FromRef(registry, repository, ref, branchPrefix, tagPrefix string) (string, error) {
	switch {
	case strings.HasPrefix(ref, "refs/heads/"):
		return fromRefBranch(registry, repository, ref, branchPrefix)
	case strings.HasPrefix(ref, "refs/tags/"):
		return fromRefTag(registry, repository, ref, tagPrefix)
	default:
		return "", errors.New("git ref is neither branch not tag")
	}
}

func fromRefBranch(registry, repository, ref, prefix string) (tag string, err error) {
	ref = strings.TrimPrefix(ref, "refs/heads/")
	tag = fmt.Sprintf("%s%s", prefix, slug.Make(ref))
	tag = prepend(registry, repository, tag)
	return
}

func fromRefTag(registry, repository, ref, prefix string) (tag string, err error) {
	ref = strings.TrimPrefix(ref, "refs/tags/")
	ref = strings.TrimPrefix(ref, prefix)
	tag = fmt.Sprintf("%s%s", prefix, ref)
	tag = prepend(registry, repository, tag)
	return
}

func prepend(registry, repository, tag string) string {
	registry = strings.TrimSuffix(registry, "/")
	repository = strings.TrimSuffix(repository, "/")
	switch {
	case registry != "" && repository != "":
		return fmt.Sprintf("%s/%s:%s", registry, repository, tag)
	case repository != "":
		return fmt.Sprintf("%s:%s", repository, tag)
	default:
		return tag
	}
}

func FromSHA(registry, repository, sha, prefix string) (tag string, err error) {
	if len(sha) != 40 {
		return "", errors.Errorf("invalid value: commit sha: %q", sha)
	}
	tag = fmt.Sprintf("%s%.7s", prefix, sha)
	tag = prepend(registry, repository, tag)
	return
}

func FromRunNumber(registry, repository string, n int, prefix string) (tag string, err error) {
	if n == 0 {
		return "", errors.Errorf("invalid value: run number: %q", n)
	}
	tag = fmt.Sprintf("%s%d", prefix, n)
	tag = prepend(registry, repository, tag)
	return
}

func output(name, value string) {
	_, _ = fmt.Fprintf(os.Stdout, "::set-output name=%s::%s\n", name, value)
}

func die(e error) {
	_, _ = fmt.Fprintf(os.Stdout, "::error ::%s\n", e.Error())
	os.Exit(-1)
}

func main() {
	var inputs struct {
		Registry   string `envconfig:"INPUT_REGISTRY"`
		Repository string `envconfig:"INPUT_REPOSITORY"`

		RefTagPrefix    string `envconfig:"INPUT_REFTAGPREFIX"`
		RefBranchPrefix string `envconfig:"INPUT_REFBRANCHPREFIX"`
		SHAPrefix       string `envconfig:"INPUT_SHAPREFIX"`
		SerialPrefix    string `envconfig:"INPUT_SERIALPREFIX"`

		GithubRef       string `envconfig:"GITHUB_REF"`
		GithubSHA       string `envconfig:"GITHUB_SHA"`
		GithubRunNumber int    `envconfig:"GITHUB_RUN_NUMBER"`
	}
	if err := envconfig.Process("", &inputs); err != nil {
		die(err)
	}

	ref, err := FromRef(inputs.Registry, inputs.Repository, inputs.GithubRef, inputs.RefBranchPrefix, inputs.RefTagPrefix)
	if err != nil {
		die(err)
	}
	sha, err := FromSHA(inputs.Registry, inputs.Repository, inputs.GithubSHA, inputs.SHAPrefix)
	if err != nil {
		die(err)
	}
	serial, err := FromRunNumber(inputs.Registry, inputs.Repository, inputs.GithubRunNumber, inputs.SerialPrefix)
	if err != nil {
		die(err)
	}

	output("ref", ref)
	output("sha", sha)
	output("serial", serial)
}
