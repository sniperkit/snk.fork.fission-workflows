
PROTO_TARGETS=$(shell find pkg -type f -name "*.pb.go")
PROTO_TARGETS+=$(shell find pkg -type f -name "*.pb.gw.go")
SRC_TARGETS=$(shell find pkg -type f -name "*.go")
CHART_FILES=$(shell find charts/fission-workflows -type f)

.PHONY: build generate prepush verify test

build wfcli fission-workflows-bundle:
	# TODO toggle between container and local build, support parameters, seperate cli and bundle
	build/build.sh

generate: ${PROTO_TARGETS} pkg/version/version.gen.go examples/workflows-env.yaml

prepush: generate verify test

test:
	test/runtests.sh

verify:
	hack/verify-workflows.sh
	hack/verify-gofmt.sh
	hack/verify-govet.sh
	helm lint charts/fission-workflows/ > /dev/null

clean:
	rm wfcli* fission-workflows-bundle*

version pkg/version/version.gen.go: pkg/version/version.go ${SRC_TARGETS}
	hack/codegen-version.sh -o pkg/version/version.gen.go -v head

examples/workflows-env.yaml: ${CHART_FILES}
	hack/codegen-helm.sh

%.swagger.json: %.pb.go
	hack/codegen-swagger.sh

%.pb.gw.go %.pb.go: %.proto
	hack/codegen-grpc.sh

# TODO add: release, docker builds, (quick) deploy, test-e2e