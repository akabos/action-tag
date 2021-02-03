package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromRef(t *testing.T) {
	t.Run("branch master", func(t *testing.T) {
		res, err := FromRef("", "", "refs/heads/master", "", "")
		require.NoError(t, err)
		require.Equal(t, "master", res)
	})
	t.Run("branch master with prefix", func(t *testing.T) {
		res, err := FromRef("", "", "refs/heads/master", "branch-", "")
		require.NoError(t, err)
		require.Equal(t, "branch-master", res)
	})
	t.Run("branch master with repository", func(t *testing.T) {
		res, err := FromRef("", "example/example", "refs/heads/master", "", "")
		require.NoError(t, err)
		require.Equal(t, "example/example:master", res)
	})
	t.Run("branch master with repository and registry", func(t *testing.T) {
		res, err := FromRef("example.com", "example/example", "refs/heads/master", "", "")
		require.NoError(t, err)
		require.Equal(t, "example.com/example/example:master", res)
	})
	t.Run("branch feature", func(t *testing.T) {
		res, err := FromRef("", "", "refs/heads/features/some-feature", "", "")
		require.NoError(t, err)
		require.Equal(t, "features-some-feature", res)
	})
	t.Run("branch release", func(t *testing.T) {
		res, err := FromRef("", "", "refs/heads/releases/1.0.0", "", "")
		require.NoError(t, err)
		require.Equal(t, "releases-1-0-0", res)
	})
	t.Run("tag", func(t *testing.T) {
		res, err := FromRef("", "", "refs/tags/1.0.0", "", "")
		require.NoError(t, err)
		require.Equal(t, "1.0.0", res)
	})
	t.Run("tag with prefix", func(t *testing.T) {
		res, err := FromRef("", "", "refs/tags/1.0.0", "", "v")
		require.NoError(t, err)
		require.Equal(t, "v1.0.0", res)
	})
	t.Run("tag with prefix collapse", func(t *testing.T) {
		res, err := FromRef("", "", "refs/tags/v1.0.0", "", "v")
		require.NoError(t, err)
		require.Equal(t, "v1.0.0", res)
	})
	t.Run("error neither branch nor tag", func(t *testing.T) {
		_, err := FromRef("", "", "refs/stash", "", "")
		require.Error(t, err)
	})
}

func TestFromSHA(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		res, err := FromSHA("", "", "8d5ef19c82f96880e47c5017afc2036ab167a9a0", "")
		require.NoError(t, err)
		require.Equal(t, "8d5ef19", res)
	})
	t.Run("with prefix", func(t *testing.T) {
		res, err := FromSHA("", "", "8d5ef19c82f96880e47c5017afc2036ab167a9a0", "sha-")
		require.NoError(t, err)
		require.Equal(t, "sha-8d5ef19", res)
	})
	t.Run("with repository", func(t *testing.T) {
		res, err := FromSHA("", "example/example", "8d5ef19c82f96880e47c5017afc2036ab167a9a0", "")
		require.NoError(t, err)
		require.Equal(t, "example/example:8d5ef19", res)
	})
	t.Run("with repository and registry", func(t *testing.T) {
		res, err := FromSHA("example.com", "example/example", "8d5ef19c82f96880e47c5017afc2036ab167a9a0", "")
		require.NoError(t, err)
		require.Equal(t, "example.com/example/example:8d5ef19", res)
	})
	t.Run("error", func(t *testing.T) {
		_, err := FromSHA("", "", "", "")
		require.Error(t, err)
	})
}

func TestFromRunNumber(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		res, err := FromRunNumber("", "", 100500, "")
		require.NoError(t, err)
		require.Equal(t, "100500", res)
	})
	t.Run("with prefix", func(t *testing.T) {
		res, err := FromRunNumber("", "", 100500, "build-")
		require.NoError(t, err)
		require.Equal(t, "build-100500", res)
	})
	t.Run("with repository", func(t *testing.T) {
		res, err := FromRunNumber("", "example/example", 100500, "")
		require.NoError(t, err)
		require.Equal(t, "example/example:100500", res)
	})
	t.Run("with repository and registry", func(t *testing.T) {
		res, err := FromRunNumber("example.com", "example/example", 100500, "")
		require.NoError(t, err)
		require.Equal(t, "example.com/example/example:100500", res)
	})
	t.Run("error", func(t *testing.T) {
		_, err := FromRunNumber("", "", 0, "")
		require.Error(t, err)
	})
}