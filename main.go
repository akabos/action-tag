package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosimple/slug"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func FromRef(ref, branchPrefix, tagPrefix string) (string, error) {
	switch {
	case strings.HasPrefix(ref, "refs/heads/"):
		return fromRefBranch(ref, branchPrefix)
	case strings.HasPrefix(ref, "refs/tags/"):
		return fromRefTag(ref, tagPrefix)
	default:
		return "", errors.New("git ref is neither branch not tag")
	}
}

func fromRefBranch(ref, prefix string) (string, error) {
	ref = strings.TrimPrefix(ref, "refs/heads/")
	return fmt.Sprintf("%s%s", prefix, slug.Make(ref)), nil
}

func fromRefTag(ref, prefix string) (string, error) {
	ref = strings.TrimPrefix(ref, "refs/tags/")
	ref = strings.TrimPrefix(ref, prefix)
	return fmt.Sprintf("%s%s", prefix, ref), nil
}

func FromSHA(sha, prefix string) (string, error) {
	if len(sha) != 40 {
		return "", errors.Errorf("invalid value: commit sha: %q", sha)
	}
	return fmt.Sprintf("%s%.7s", prefix, sha), nil
}

func FromRunNumber(n int, prefix string) (string, error) {
	if n == 0 {
		return "", errors.Errorf("invalid value: run number: %q", n)
	}
	return fmt.Sprintf("%s%d", prefix, n), nil
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

	ref, err := FromRef(inputs.GithubRef, inputs.RefBranchPrefix, inputs.RefTagPrefix)
	if err != nil {
		die(err)
	}
	sha, err := FromSHA(inputs.GithubSHA, inputs.SHAPrefix)
	if err != nil {
		die(err)
	}
	serial, err := FromRunNumber(inputs.GithubRunNumber, inputs.SerialPrefix)
	if err != nil {
		die(err)
	}

	output("ref", ref)
	output("sha", sha)
	output("serial", serial)
}
