#!/bin/sh

echo '[#] compiling in go'
protoc --proto_path=protobuf --go_out=back/proto --go_opt=paths=source_relative protobuf/*

echo '[#] compiling in ts'
protoc --proto_path=protobuf --plugin=./front/node_modules/.bin/protoc-gen-ts_proto --ts_proto_out front/src/proto protobuf/*

echo '[+] done!'