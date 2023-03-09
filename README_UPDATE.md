<!-- markdownlint-configure-file {
  "MD013": {
    "code_blocks": false,
    "tables": false
  },
  "MD033": false,
  "MD041": false
} -->

<div align="center">

# Carpenter

Carpenter is a tool to automatically test, build, and publish Arcaflow plugins.
 
</div>

## Table of Contents

[Requirements](#requirements) •
[Configuration](#configuration) •
[Build Carpenter Executable](#build-carpenter-executable) •
[Build Carpenter Image](#build-carpenter-image) •
[Example Build Execution With Executable](#example-build-execution-with-executable)
[Example Build and Push Execution Containerized](#example-build-and-push-execution-containerized)
[Using Carpenter As Reusable Workflow](#using-carpenter-in-reusable-workflow) •

## Requirements

* golang v1.18
* docker
* python 3 and pip
* flake8
* current working directory is this project's root directory

Each plugin directory must meet the [Arcaflow Plugins Requirements](https://github.com/arcalot/arcaflow-plugins#requirements-for-plugins).

## Configuration

Configuring the carpenter can be done in the `carpenter.yaml` file

example `.carpenter.yaml`
```yaml
revision: 20220824
image_name: arcaflow-plugin-template-python
image_tag: '0.0.1'
project_filepath: ../arcaflow-plugin-template-python
registries:
  - url: ghcr.io
    username_envvar: "GITHUB_USERNAME"
    password_envvar: "GITHUB_PASSWORD"
  - url: quay.io
    username_envvar: "QUAY_USERNAME"
    password_envvar: "QUAY_PASSWORD"
    namespace_envvar: "QUAY_NAMESPACE"
```
### Environment variables

Carpenter additionally can be configured beyond defaults using enviornment variables

  - Configureable Environment Variables:
    | Env Var               | Description                                 | Type   | Default Value                |
    | --------------------- | ------------------------------------------- | ------ | ---------------------------- |
    | IMAGE_TAG             | Tag of image being built                    | String | image_tag in carpenter.yaml  |
    | IMAGE_NAME            | Name of image being built                   | String | image_name in carpenter.yaml |
    | BUILD_TIMEOUT         | Length of time before build fails           | Int    | 600                          |
    | QUAY_IMG_EXP          | Image label to automatically expire in Quay | String | "never"                      |
    | QUAY_CUSTOM_NAMESPACE | Overrides $QUAY_NAMESPACE                   | String | ""                           |

* `BUILD_TIMEOUT` accepts an integer representing the number of **seconds** before a build timeouts.
* `QUAY_IMG_EXP` more documentation and time formats can be found [here](https://docs.projectquay.io/use_quay.html#:~:text=Setting%20tag%20expiration%20from%20a%20Dockerfile)
* `QUAY_CUSTOM_NAMESPACE` if set, will use in place of `QUAY_NAMESPACE`. More info [Using Carpenter As Reusable Workflow](#using-carpenter-in-reusable-workflow) •

## Build Carpenter Executable

```shell
go build carpenter.go
```

If successful, this will result in the carpenter executable named `carpenter` in your current working directory.

## Build Carpenter Image

```shell
docker build . --tag carpenter-img
```

## Example Build Execution With Executable

example `.carpenter.yaml`
```yaml
revision: 20220824
image_name: arcaflow-plugin-template-python
image_tag: '0.0.1'
project_filepath: ../arcaflow-plugin-template-python
registries:
  - url: ghcr.io
    username_envvar: "GITHUB_USERNAME"
    password_envvar: "GITHUB_PASSWORD"
  - url: quay.io
    username_envvar: "QUAY_USERNAME"
    password_envvar: "QUAY_PASSWORD"
    namespace_envvar: "QUAY_NAMESPACE"
```

```shell
./carpenter build --build
```

## Example Build and Push Execution Containerized

```shell
docker run \
    --rm \
    -e=IMAGE_TAG="0.0.1"\
    -e=BUILD_TIMEOUT=600\
    -e=GITHUB_USERNAME=$GITHUB_USERNAME \
    -e=GITHUB_PASSWORD=$GITHUB_PASSWORD \
    -e=QUAY_USERNAME=$QUAY_USERNAME\
    -e=QUAY_PASSWORD=$QUAY_PASSWORD\
    -e=QUAY_NAMESPACE=$QUAY_NAMESPACE\
    --volume /var/run/docker.sock:/var/run/docker.sock:z \
    --volume $PWD/../arcaflow-plugin-template-python:/github/workspace \
    carpenter-img build --build --push
```

## Using Carpenter In Reusable Workflow

From within a plugin repository you can utilize carpenter to test, build, and push automatically.
Secrets should be configuerd within the repository for credentials.

```yaml
name: Carpenter
on:
  push:
    branches:
      - "**"
  release:
    types:
      - published

jobs:
  carpenter:
    uses: arcalot/carpenter/.github/workflows/carpenter_reusable_workflow.yaml@main
    with:
      image_name: ${{ github.event.repository.name }}
      image_tag: 'latest'
      quay_img_exp: 'never'
      github_username: ${{ github.actor }}
      github_namespace: ${{ github.repository_owner }}
      quay_custom_namespace: 'example' # This is optional, for reference
    secrets: 
      QUAY_NAMESPACE: ${{ secrets.QUAY_NAMESPACE }}
      QUAY_USERNAME: ${{ secrets.QUAY_USERNAME }}
      QUAY_PASSWORD: ${{ secrets.QUAY_PASSWORD }}

```
