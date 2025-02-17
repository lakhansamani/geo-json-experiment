all: build
build:
	protoc \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative \
  --plugin=./node_modules/.bin/protoc-gen-ts_proto \
  --ts_proto_out=./typescript \
  --ts_proto_opt=esModuleInterop=true,outputServices=grpc-js,useExactTypes=false \
  apis/layer/layer.proto
start_db:
	@docker run --name test-db -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres -d postgis/postgis
stop_db:
	@docker stop test-db
	@docker rm test-db
run:
	go run main.go