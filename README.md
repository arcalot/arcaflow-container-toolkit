# Carpenter, Arcaflow Plugin Image Builder

Arcaflow Plugin Image Builder is a tool which has been developed for automatically testing, building, and publishing Arcaflow plugins.

More in detail:
* Python plugins are going to be unit tested, scanned with pyflakes, and coverage data for each plugin will be collected.
  Successfully tested plugins will be published on pypi registry automatically on tag event.
* Go plugins are going to be unit tested and coverage data for each plugin will be collected.

Successfully tested plugins will be also added to docker images and end to end tested where possible.
Successfully tested images will be published to quay.io automatically on tag event.

# Preparing the project for being built with Arcaflow Plugin Image Builder

Each plugin directory must meet the [Arcaflow Plugins Requirements](https://github.com/arcalot/arcaflow-plugins#requirements-for-plugins).

The builder will check that the requirements are met.

## Test carpenter

* golang v1.18
* current working directory is this project's root directory
* mock interfaces for carpenter's interfaces
* python 3 and pip

Install flake8
```shell
python3 -m pip install --user flake8
```

Generate golang mocks for carpenter's interfaces
```shell
go generate ./...
```

Execute test suite with statement coverage and save coverage data to `cov.txt`
```shell
go test ./... -coverprofile=cov.txt
```

## Build the carpenter

* golang v1.18
* current working directory is this project's root directory
* flake8

```shell
go build carpenter.go
```

If successful, this will result in the arcaflow-plugin-image-builder executable, and it will be named `carpenter` in your current working directory.

## Build carpenter's image

### Requirements

* docker
* current working directory is `arcaflow-plugin-image-builder`

```shell
docker build . --tag carpenter-img
```

## Example Build Execution

### Requirements

* arcaflow-plugin-image-builder executable named `carpenter`
* `.carpenter.yaml` in the same directory as `carpenter` executable
* flake8

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

### Requirements

* docker engine and cli
* a carpenter image named `carpenter-img`
* plugin project
* GitHub username and password
* Quay username and password

```shell
docker run \
    --rm \
    -e=IMAGE_TAG="0.0.1"\
    -e=GITHUB_USERNAME=$GITHUB_USERNAME \
    -e=GITHUB_PASSWORD=$GITHUB_PASSWORD \
    -e=QUAY_USERNAME=$QUAY_USERNAME\
    -e=QUAY_PASSWORD=$QUAY_PASSWORD\
    -e=QUAY_NAMESPACE=$QUAY_NAMESPACE\
    --volume /var/run/docker.sock:/var/run/docker.sock:z \
    --volume $PWD/../arcaflow-plugin-template-python:/github/workspace \
    carpenter-img build --build --push
```

You can override the `image_tag` from `.carpenter.yaml` by injecting the environment variable `IMAGE_TAG`, set to your chosen string, into
`carpenter-img` when you run the container. The same goes for `image_name` and `IMAGE_NAME`, and your registry environment variables.
