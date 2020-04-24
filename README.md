# action-tag

GitHub action which generates Docker image tags from Git commit metadata.

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
