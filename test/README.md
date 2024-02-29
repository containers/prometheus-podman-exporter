# exporter e2e tests with bats

## Running tests

To run the tests locally in your sandbox first build and start the exporter:

```shell
$ make binary
$ ./bin/prometheus-podman-exporter --collector.enable-all --collector.store_labels --debug
```

After starting the exporter run:

```shell
$ make test-e2e
```
