.DELETE_ON_ERROR:

target := headless-askpass
srcs := $(wildcard *.go)
arches := amd64 arm arm64

default: $(target).$(shell go env GOOS)-$(shell go env GOARCH)

.PHONY: cross
cross: $(addprefix $(target).linux-,$(arches))

$(target).%: $(srcs)
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) \
	go build -o $@
