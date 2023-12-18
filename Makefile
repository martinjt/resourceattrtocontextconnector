IMAGE_REPOSITORY ?= "otelcol-custom"
IMAGE_TAG ?= "latest"

.PHONY: build
build:
	builder --config=ocb.yaml

.PHONY: build-and-push-dockerhub
build-and-push-dockerhub: build
	docker build -t $(IMAGE_REPOSITORY):$(IMAGE_TAG) .
	docker push $(IMAGE_REPOSITORY):$(IMAGE_TAG)

# This is required to build the a proper executable that can be included in a docker image for an arm64 architecture K8s cluster
.PHONY: build-arm
build-arm:
	env GOOS=linux GOARCH=arm64 builder --config=ocb.yaml

.PHONY: build-arm-and-push-dockerhub
build-arm-and-push-dockerhub: build-arm
	docker buildx build --push --platform linux/arm/v7,linux/arm64/v8 -t $(IMAGE_REPOSITORY):$(IMAGE_TAG) .

.PHONY: generate-metadata
generate-metadata:
	cd resourceattrtocontextconnector && go generate

.PHONY: run
run:
	./dist/otelcol-custom --config=test-config.yaml

.PHONY: debug
debug:
	dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient --log exec ./dist/otelcol-custom -- --config=config.yaml

.PHONY: test-span
test-span:
	otel-cli span --endpoint localhost:4317 --service "test" --name "Test Span" --attrs "app.honeycomb_api_key=$(HONEYCOMB_API_KEY)" --tp-print