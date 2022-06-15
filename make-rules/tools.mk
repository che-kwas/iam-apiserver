# Makefile helper functions for tools

TOOLS := golangci-lint mockgen codegen protoc protoc-gen-go protoc-gen-go-grpc

.PHONY: tools.verify
tools.verify: $(addprefix tools.verify., $(TOOLS))

.PHONY: tools.verify.%
tools.verify.%:
	@type $* >/dev/null 2>&1 || make tools.install.$*

.PHONY: tools.install.%
tools.install.%:
	make install.$*

.PHONY: install.golangci-lint
install.golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: install.mockgen
install.mockgen:
	go install github.com/golang/mock/mockgen@latest

.PHONY: install.codegen
install.codegen:
	go install github.com/che-kwas/iam-kit/tools/codegen@latest

.PHONY: install.protoc
install.protoc:
	sudo apt-get -y install protobuf-compiler

.PHONY: install.protoc-gen-go
install.protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

.PHONY: install.protoc-gen-go-grpc
install.protoc-gen-go-grpc:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
