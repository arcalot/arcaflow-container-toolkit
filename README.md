# Carpenter, Arcaflow Plugin Image Builder

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