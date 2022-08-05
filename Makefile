nats_storage = stan/store
nats_cluster = order-cluster
nats_max_msgs = 1000
nats_max_bytes = 100000k

.PHONY: build
build:
	go build -v ./cmd/apiserver  &&   mv apiserver ./bin/apiserver

.DEFAULT_GOAL := build

up:
	./bin/apiserver

.PHONY: buildpub
buildpub:
	go build -v ./cmd/publisher  &&   mv publisher ./bin/publisher

uppub:
	./bin/publisher

.PHONY: stan
stan:
	nats-streaming-server -cid $(nats_cluster) -store file -dir $(nats_storage) -max_msgs $(nats_max_msgs) -max_bytes $(nats_max_bytes)

su:
	make build && make up

pu:
	make buildpub && make uppub

test:
	go test -v ./...

.PHONY: postattack
postattack:
	go run ./vegeta/post/attack.go

.PHONY: getattack
getattack:
	go run ./vegeta/get/attack.go