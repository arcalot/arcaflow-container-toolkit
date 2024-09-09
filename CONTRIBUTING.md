# Contributing

## Generating Mocks

As of writing, the mock objects are generated in this project's `mock` directory. When a given component is altered and its test suite depends upon a mock object new mocks need to be generated in the `mock` directory to replace the now outdated mocks. 

To generate new mocks, while in this project's root directory execute this command:

```shell
go generate ./...
```
