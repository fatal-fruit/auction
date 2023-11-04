# Auction Module

### Developing

Generate mocks
```shell
mockgen --source=./types/expected_keepers.go --destination=./testutil/expected_keepers_mocks.go --package=testutil
mockgen --source=./types/expected_services.go --destination=./testutil/expected_services_mocks.go --package=testutil
```