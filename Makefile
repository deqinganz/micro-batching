.PHONY: codegen
codegen:
	oapi-codegen -generate types,gin,spec,strict-server -package api api/openapi.yaml > api/api.gen.go

.PHONY: clean
clean:
	rm -rf api/api.gen.go