# drone-docker

[![Build Status](https://github.com/kit101/drone-plugin-docker-cache/actions/workflows/publish.yml/badge.svg)](https://github.com/kit101/drone-plugin-docker-cache/actions/workflows/publish.yml)
[![](https://images.microbadger.com/badges/image/plugins/docker.svg)](https://microbadger.com/images/plugins/docker "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/kit101/drone-plugin-docker-cache?status.svg)](http://godoc.org/github.com/kit101/drone-plugin-docker-cache)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-docker)](https://goreportcard.com/report/github.com/kit101/drone-plugin-docker-cache)

用于将已有docker data root缓存复制到registry path并处理.dockerignore


## Docker

Build the Docker images with the following [docker/Dockerfile](https://github.com/kit101/drone-plugin-docker-cache/blob/master/docker/Dockerfile)
```shell
# docker buildx
docker buildx build --platform linux/amd64,linux/arm64 -t kit101z/drone-plugin-docker-cache -f docker/Dockerfile . --push
# docker build
docker build --platform linux/amd64,linux/arm64 -t kit101z/drone-plugin-docker-cache -f docker/Dockerfile . --push
```
## Usage

### Using in drone

```yaml
kind: pipeline
name: default

steps:
- name: build dummy docker file and publish
  image: kit101z/drone-plugin-docker-cache
  pull: never
  settings:
    registry-path: dockerlib
    src: /mnt/dockerlib
  volumes:
  - name: dockerlib
    path: /mnt/dockerlib
volumes:
- name: dockerlib
  host:
    path: /var/drone-runner/cache/dockerlib

```
### Using in docker
```shell
➜  ~ docker run --rm \
        -v $(PWD)/.dockerlibs/:/mnt/dockerlib \
        -w /drone/src \
        -e PLUGIN_STORAGE_PATH=dockerlib \
        -e PLUGIN_SRC=/mnt/dockerlib kit101z/drone-docker-cache
```

### Usage in cmd
```shell
➜  drone-plugins-docker-cache git:(main) ✗ release/darwin/amd64/docker-cache --help                                          
NAME:
   docker cache plugin - docker data root缓存插件，复制已有docker data root目录到registry path并处理.dockerignore

USAGE:
   docker-cache [global options] command [command options] [arguments...]

VERSION:
   unknown

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --storage-path value  docker data root 目录 [$PLUGIN_STORAGE_PATH]
   --src value           docker data root缓存目录.
该值存在会复制缓存目录到storage-path. 
若storage-path是在workingDir下，则还会在的${workingDir}/.dockerignore中追加storage-path的相对路径. [$PLUGIN_SRC]
   --dockerignores value  .dockerignore中额外写入的忽略路径 [$PLUGIN_DOCKERIGNORES]
   --help, -h             show help
   --version, -v          print the version

```

