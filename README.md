<p align="center"><img width=50% src="img/act-logo.png"></p>

<!-- markdownlint-configure-file {
  "MD013": {
    "code_blocks": false,
    "tables": false
  },
  "MD033": false,
  "MD041": false
} -->

<div align="center">

Arcaflow Container Toolkit is a tool to automatically test, build, and publish Arcaflow plugins to a single or multiple repositories.
 
</div>

## Table of Contents

• [Requirements](#requirements)  
• [Configuration](#configuration)  
• [Build Arcaflow Container Toolkit As Executable Locally](#build-arcaflow-container-toolkit-as-executable-locally)  
• [Arcaflow Container Toolkit as a Package](#arcaflow-container-toolkit-as-a-package)  
• [Arcaflow Container Toolkit and Reusable Workflows](#arcaflow-container-toolkit-and-reusable-workflows)  
• [Arcaflow Container Toolkit as an Action](#arcaflow-container-toolkit-as-an-action)  

## Requirements

* golang v1.18
* docker
* python 3 and pip
* flake8
* current working directory is this project's root directory for local builds

Each plugin directory must meet the [Arcaflow Plugins Requirements](https://github.com/arcalot/arcaflow-plugins#requirements-for-plugins).

## Configuration

Configuring Arcaflow Container Toolkit can be done in the `act.yaml` file as well as setting environment variables.

### Configurable Variables

#### Required:
  `IMAGE_NAME` Name of the image that Arcaflow Container Toolkit will build - string  
  `IMAGE_TAG`  Tag of the image that Arcaflow Container Toolkit will build - string  
#### Optional:  
  `GITHUB_USERNAME` Github Username to be used for credentials - Default: ""  
  `GITHUB_PASSWORD` Github Password to be used for credentials - Default: ""  
  `GITHUB_NAMESPACE` Github Namespace to push image - Default: ""  
  `QUAY_USERNAME` Quay Username to be used for credentials - Default: ""  
  `QUAY_PASSWORD` Quay Password to be used for credentials - Default: ""  
  `QUAY_NAMESPACE` Quay Namespace to push image - Default: ""  
  `QUAY_CUSTOM_NAMESPACE` Quay Namespace to push image that is not QUAY_NAMESPACE - Default: ""  
  `QUAY_IMG_EXP` Image label to automatically expire in Quay - Default: "never"  
  `BUILD_TIMEOUT` Length of time before a build will fail in seconds - Default: 600  

#### Additional Information
* `QUAY_IMG_EXP` more documentation and time formats can be found [here](https://docs.projectquay.io/use_quay.html#:~:text=Setting%20tag%20expiration%20from%20a%20Dockerfile)
* `QUAY_CUSTOM_NAMESPACE` if set, will use in place of `QUAY_NAMESPACE`. More info [Arcaflow Container Toolkit and Reusable Workflows](#arcaflow-container-toolkit-and-reusable-workflows)

## Build Arcaflow Container Toolkit As Executable Locally

Arcaflow Container Toolkit can be ran locally by building an executable.  
Configure `act.yaml` and or set environment variables.  

```shell
vi act.yaml
```

example `.act.yaml`
```yaml
revision: 20220824
image_name: "<IMAGE_NAME>"
image_tag: "<IMAGE_TAG>"
project_filepath: "<path/to/plugin/project/>"
registries:
  - url: ghcr.io
    username_envvar: "<GITHUB_USERNAME>"
    password_envvar: "<GITHUB_PASSWORD>"
  - url: quay.io
    username_envvar: "<QUAY_USERNAME>"
    password_envvar: "<QUAY_PASSWORD>"
    namespace_envvar: "<QUAY_NAMESPACE>"
```

#### Build the executable

```shell
go build act.go
```
#### Arcaflow Container Toolkit test and build

```shell
./act build --build
```

#### Arcaflow Container Toolkit test, build, and push

```shell
./act build --build --push
```

## Arcaflow Container Toolkit as a Package

Pull the latest image

```shell
docker pull ghcr.io/arcalot/arcaflow-container-toolkit:latest
```

Run the Arcaflow Container Toolkit image with enviornment variables

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
    ghcr.io/arcalot/arcaflow-container-toolkit:latest build --build --push
```

## Arcaflow Container Toolkit and Reusable Workflows

Arcaflow Container Toolkit can be utilized using the official reusable workflow `arcalot/arcaflow-containter-toolkit/.github/workflows/reusable_workflow.yaml`.  

```yaml
name: Arcaflow Container Toolkit
on:
  push:
    branches:
      - "**"
  release:
    types:
      - published

jobs:
  arcaflow-container-toolkit:
    uses: arcalot/arcaflow-container-toolkit/.github/workflows/reusable_workflow.yaml@main
    with:
      image_name: ${{ github.event.repository.name }}
      image_tag: 'latest'     
    secrets: 
      QUAY_NAMESPACE: ${{ secrets.QUAY_NAMESPACE }}
      QUAY_USERNAME: ${{ secrets.QUAY_USERNAME }}
      QUAY_PASSWORD: ${{ secrets.QUAY_PASSWORD }}

```

## Arcaflow Container Toolkit as an Action

Arcaflow Container Toolkit can be utilized as an action in a workflow.  

```yaml
- name: arcaflow-container-toolkit-action
        uses: arcalot/arcaflow-container-toolkit-action@v0.1.0
        with:
          image_name: ${{ github.event.repository.name }}
          image_tag: 'latest'
          quay_username: ${{ secrets.QUAY_USERNAME }}
          quay_password: ${{ secrets.QUAY_PASSWORD }}
          quay_namespace: ${{ secrets.QUAY_NAMESPACE }}
          quay_custom_namespace: 'example' # This is optional, for reference

```
