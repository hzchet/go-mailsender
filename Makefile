.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml

.PHONY: formate
formate:
	gofumpt -l -w .

.PHONY: swagger
swagger:
	echo "swag gen"