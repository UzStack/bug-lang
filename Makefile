build:
	go build -o buglang ./cmd/main.go

air:
	air --build.cmd="go build -o /tmp/buglang ./cmd/main.go" --build.bin="/tmp/buglang" .

air-fpm:
	air --build.cmd="go build -o /tmp/buglang ./cmd/fpm/main.go" --build.bin="/tmp/buglang" .

.PHONY: build