# Run tests with coverage
test:
    go test -cover -bench=. -benchmem -race ./... -coverprofile=coverage.out

# Build treekanga binary to GOPATH/bin
build:
    go build -buildvcs=false -ldflags "-X 'main.version=`git describe --tags --abbrev=0`'" -o `go env GOPATH`/bin/treekanga
