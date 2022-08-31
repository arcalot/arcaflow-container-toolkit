# carpenter

## Build the builder

* golang v1.17
* current working directory is this project's root directory

```shell
go build build.go
```

## Example build configurations

### Build a single plugin directory

example `build.yaml`
```yaml
revision: 20220824
target: fio
project_filepath: ../arcaflow-plugins/python/fio
```

### Build a directory of plugins

example `build.yaml`
```yaml
revision: 20220824
target: all
project_filepath: ../arcaflow-plugins/python
```

## Example execution

### Requirements

* `build.yaml` in the same directory with this executable

```shell
./build
```