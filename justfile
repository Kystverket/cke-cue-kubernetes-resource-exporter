[private]
default:
  @just --list -u

# Output test cue to stdout
run:
  go run cmd/cke/go

# Output test cue to _rendered/
render:
  go run cmd/cke/cke.go -out files

# build and install go binary
build:
  go build cmd/cke/cke.go

# Install helper
install:
  go install cmd/cke/cke.go
