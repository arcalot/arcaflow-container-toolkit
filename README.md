# carpenter

## Build the builder

* golang v1.17
* current working directory is this project's root directory

```shell
go build build.go
```

## Example Execution

### Requirements

* `build.yaml` in the same directory with this executable

example `build.yaml`
```yaml
revision: 20220824
target: fio
project_filepath: ../arcaflow-plugins/python/fio
```

```shell
./build
```