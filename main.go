package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvPrefix("ACTION_TAG")
}

const (
	optRefBranchPrefix = "ref_branch_prefix"
	optRefTagPrefix    = "ref_tag_prefix"
	optSHAPrefix       = "sha_prefix"
	optSerialPrefix    = "serial_prefix"
	optGithubSHA       = "github_sha"
	optGithubRef       = "github_ref"
	optGithubRunNumber = "github_run_number"
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

func exit(e error) {
	_, _ = fmt.Fprintf(os.Stdout, "::error ::%s\n", e.Error())
	os.Exit(-1)
}

func main() {
	ref, err := FromRef(viper.GetString(optGithubRef), viper.GetString(optRefBranchPrefix), viper.GetString(optRefTagPrefix))
	if err != nil {
		exit(err)
	}
	sha, err := FromSHA(viper.GetString(optGithubSHA), viper.GetString(optSHAPrefix))
	if err != nil {
		exit(err)
	}
	serial, err := FromRunNumber(viper.GetInt(optGithubRunNumber), viper.GetString(optSerialPrefix))
	if err != nil {
		exit(err)
	}
	output("ref", ref)
	output("sha", sha)
	output("serial", serial)
}
