# Carpenter, Arcaflow Plugin Image Builder

## Build Carpenter Image

### Requirements

* docker
* current working directory is `arcaflow-plugin-image-builder`

```shell
docker build . --tag carpenter-img
```

## Example build configurations

### Build a single plugin directory

example `.carpenter.yaml`
```yaml
revision: 20220824
image_name: arcaflow-plugin-template-python
image_tag: '0.1.0'
project_filepath: ../arcaflow-plugin-template-python
registries:
  - url: ghcr.io
    username_envvar: "GITHUB_USERNAME"
    password_envvar: "GITHUB_PASSWORD"
  - url: quay.io
    username_envvar: "QUAY_USERNAME"
    password_envvar: "QUAY_PASSWORD"
```

## Example Execution Containerized

### Requirements

* docker engine and cli
* a carpenter image named `carpenter-img`
* plugin project
* GitHub username and password
* Quay username and password

```shell
docker run \
    --rm \
    -e=IMAGE_TAG="0.1.1"\
    -e=GITHUB_USERNAME=$GITHUB_USERNAME \
    -e=GITHUB_PASSWORD=$GITHUB_PASSWORD \
    -e=QUAY_USERNAME=$QUAY_USERNAME\
    -e=QUAY_PASSWORD=$QUAY_PASSWORD\
    --volume /var/run/docker.sock:/var/run/docker.sock:z \
    --volume $PWD/../arcaflow-plugin-template-python:/github/workspace:z \
    carpenter-img build --build --push
```

You can override the `image_tag` from `.carpenter.yaml` by injecting the
environment variable `IMAGE_TAG`, set to your chosen string, into
`carpenter-img` when you run the container. The same goes for `image_name`
and `IMAGE_NAME`.