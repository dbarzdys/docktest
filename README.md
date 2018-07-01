# docktest

Use docker to run integration tests againts other services.

**Currently in beta!**

## Install

    $ go get -u github.com/dbarzdys/docktest
 

## Commands

    $ docktest
    Use docker to run integration tests againts other services.

    Usage:
    docktest [command]

    Available Commands:
    exec        Spin-up test containers and export variables to next command
    help        Help about any command
    rm          Removes containers created by DockTest
    up          Spin up test containers and export variables to .env file
    version     Show DockTest version information

    Flags:
    -c, --config string   config file (default "./docktest.yaml")
    -h, --help            help for docktest

    Use "docktest [command] --help" for more information about a command.


## Configuration

Configuration file example. Default file name: docktest.yml

``` yaml
# Extend config from other file
extend: ./docktest.other.yml
# Constants
constants:
  DB_USER: postgres
  DB_PASS: postgres
  DB_NAME: postgres
  DB_PORT: 5432
  ZIPKIN_PORT: 9411
  AUTH_PORT: 443
# Export environment variables for test
export:
  ZIPKIN_HOST: ${services.zipkin.ip}
  ZIPKIN_PORT: ${constants.ZIPKIN_PORT}
  DB_DRIVER: postgres
  DB_HOST: ${services.postgres.ip}
  DB_PORT: ${constants.DB_PORT}
  DB_USER: ${constants.DB_USER}
  DB_PASS: ${constants.DB_PASS}
  DB_NAME: ${constants.DB_NAME}
# Services used for tests
services:
  # Database
  postgres:
    image: awpc/postgres
    tag: 9.6
    env:
      POSTGRES_USER: ${constants.DB_USER}
      POSTGRES_PASSWORD: ${constants.DB_PASS}
      POSTGRES_DB: ${constants.DB_NAME}
  # Tracing
  zipkin:
    image: openzipkin/zipkin
    tag: latest
    env:
      STORAGE_TYPE: mem

```
