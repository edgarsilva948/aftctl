#
# Copyright © 2023 Edgar Costa edgarsilva948@gmail.com


# Ensure go modules are enabled:
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org

# Disable CGO so that we always generate static binaries:
export CGO_ENABLED=0

# Unset GOFLAG for CI and ensure we've got nothing accidently set
unexport GOFLAGS

aftctl:
	go build ./cmd/aftctl

fmt:
	gofmt -s -l -w cmd pkg

test:
	go test ./...		