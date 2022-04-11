# squad-dashboard

Monorepo inspiration: https://github.com/flowerinthenight/golang-monorepo

## Updating BUILD files

````
 # bazel run //:gazelle

````

## Building and Running

Building all services

````
 # bazel build //...
````

Building an individual service

````
 # bazel build //cmd/jiracollector
````

Running a service (locally)

````
 # bazel run //cmd/jiracollector
````

Building a docker container (of a service, to a tar)

````
 # bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/jiracollector:image
````

Loading a docker container in local docker registry

````
 # bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/jiracollector:image
````

Running tests

````
 # bazel test //pkg/util:util_test --test_output=errors
````

## Running all services

1. build and run all the containers (see above)
2. run docker compose to start the db
3. run docker compose to run the other services

````
 # docker compose up -d db
 # docker compose up
````
