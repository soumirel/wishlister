#!/bin/sh
set -e

PROTOC="protoc"

"$PROTOC --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go --go-grpc_opt=paths=source_relative wishlist.proto"

echo "Done"
