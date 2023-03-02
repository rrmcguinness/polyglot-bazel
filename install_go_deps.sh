#!/bin/sh

# This assumes go is on you path
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/bazelbuild/buildtools/buildifier@latest
go install github.com/bazelbuild/buildtools/buildozer@latest
go install github.com/bazelbuild/buildtools/unused_deps@latest