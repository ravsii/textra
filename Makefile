lint:
	golangci-lint run ./... -c ./.golangci.yml -v --issues-exit-code=0 --out-format colored-line-number
lint-docker:
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.52.2 golangci-lint run ./... -v --issues-exit-code=0 --out-format colored-line-number