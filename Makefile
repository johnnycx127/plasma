deps:
	@$(MAKE) -C ./contracts deps
	@echo "--> Installing Go dependencies..."
	@dep ensure -v

migrate:
	$(MAKE) -C ./contracts migrate

build:
	go build -o ./build/plasma ./cmd/plasma/main.go

abigen:
	$(MAKE) -C ./contracts abigen

protogen:
	protoc -I rpc/proto rpc/proto/root.proto --go_out=plugins=grpc:rpc/pb

build-all: abigen build

start: compile
	@./bin/start

clean:
	$(MAKE) -C ./contracts clean
	rm -rf ~/.plasma

fresh-start: clean start

.PHONY: build deps