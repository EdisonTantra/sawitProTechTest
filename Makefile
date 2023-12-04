

.PHONY: clean all init generate generate_mocks

all: build/main

build/main: main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: generate
	go mod tidy
	go mod vendor

test:
	go test -short -coverprofile coverage.out -v ./...

generate: generated generate_mocks generate_certs

generated: api.yml
	@echo "Generating files..."
	mkdir -p generated/cert || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

INTERFACES_GO_FILES := $(shell find core/port -name "port.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

generate_certs:
	openssl genrsa -out generated/cert/sawitapp 4096
	openssl rsa -in generated/cert/sawitapp -pubout -out generated/cert/sawitapp.pub
