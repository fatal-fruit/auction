## Developing for `x/action`

### Generating Protos
Protos are defined in [`/proto`](./proto)

Run the protogen make target to generate proto & pulsar files.
```shell
make proto-gen
```

### Generating Mocks
Any expected keepers must be defined in [`/types/expected_keepers.go`](./types/expected_keepers.go).
To generate run:
```
mockgen --source=./types/expected_keepers.go --destination=./testutil/expected_keepers_mocks.go --package=testutil
```