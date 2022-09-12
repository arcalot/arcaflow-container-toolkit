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

## Build the carpenter

* golang v1.17
* current working directory is this project's root directory

```shell
go build carpenter.go
```

If successful, this will result in the arcaflow-plugin-image-builder executable, and it will be named `build` in your current working directory.

## Build carpenter's image dependencies

### Requirements

* docker or podman
* current working directory is `arcaflow-plugin-image-builder/containerfiles/python`

```shell
docker build . --tag build-py
```

## Example build configurations

### Build a single plugin directory

example `.carpenter.yaml`
```yaml
revision: 20220824
target: fio
project_filepath: ../arcaflow-plugins/python/fio
```

### Build a directory of plugins

example `.carpenter.yaml`
```yaml
revision: 20220824
target: all
project_filepath: ../arcaflow-plugins/python
```

## Example execution

### Requirements

* `.carpenter.yaml` in the same directory with the `arcaflow-plugin-image-builder` executable
* `docker`, or `alias docker=podman`, executable on your `$PATH`
* a [directory of sub-directories where each sub-directory contains a Dockerfile](https://github.com/arcalot/arcaflow-plugins/tree/main/python)

```shell
./carpenter build --config .carpenter.yaml
```
