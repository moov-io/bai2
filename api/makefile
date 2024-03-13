.PHONY: client
client:
ifeq ($(OS),Windows_NT)
	@echo "Please generate ../pkg/client/ on macOS or Linux, currently unsupported on windows."
else
# Versions from https://github.com/OpenAPITools/openapi-generator/releases
#  npm install -g @openapitools/openapi-generator-cli
	@rm -rf ../pkg/client/
	OPENAPI_GENERATOR_VERSION=7.4.0 openapi-generator-cli generate --package-name client -i ./api.yaml -g go -o ../pkg/client/
	rm -rf ../pkg/client/go.mod ../pkg/client/go.sum ../pkg/client/api/ ../pkg/client/.travis.yml
	go fmt ../...
endif
