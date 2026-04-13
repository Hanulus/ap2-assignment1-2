#!/bin/bash
# Run this script once to generate Go code from .proto files.
# Requirements: protoc, protoc-gen-go, protoc-gen-go-grpc
#
# Install plugins if not present:
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

set -e

PROTO_DIR="./proto-repo"
OUT_DIR="./generated-repo"

echo "Generating payment.proto..."
protoc \
  --go_out="$OUT_DIR/payment" --go_opt=paths=source_relative \
  --go-grpc_out="$OUT_DIR/payment" --go-grpc_opt=paths=source_relative \
  -I "$PROTO_DIR" \
  "$PROTO_DIR/payment.proto"

echo "Generating orderstream.proto..."
protoc \
  --go_out="$OUT_DIR/orderstream" --go_opt=paths=source_relative \
  --go-grpc_out="$OUT_DIR/orderstream" --go-grpc_opt=paths=source_relative \
  -I "$PROTO_DIR" \
  "$PROTO_DIR/orderstream.proto"

echo "Done! Generated files are in $OUT_DIR"
