# action-tag

GitHub action which generates Docker image tags from Git commit metadata.

## Motivation

This action basically does what [docker/build-push-action](https://github.com/docker/build-push-action) provides as 
built-in functionality. So why separate action? Unfortunately, [docker/build-push-action](https://github.com/docker/build-push-action)
doesn't provide any output, so you can't use generated tags in further workflow steps. Workaround could be assigning 
a commit hash, or a build number as an additional tag, but these tend to be less informative than semantic versions.

## Usage

```yaml
- name: Generate tags
  id: tags
  uses: akabos/action-tag@master

- name: Build
  run: |
    docker build . \
      -t ${GITHUB_REPOSITORY}:${{ steps.tags.outputs.ref }} \ 
      -t ${GITHUB_REPOSITORY}:${{ steps.tags.outputs.sha }} \ 
      -t ${GITHUB_REPOSITORY}:${{ steps.tags.outputs.serial }}
```
