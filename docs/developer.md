## Message passing

### dependencies

Herald uses [Google Protocol Buffers v3](https://developers.google.com/protocol-buffers) for serialisation and [gRPC](https://grpc.io/) for message passing. Current versions:

- Protobuf 3.6.1
- gRPC 1.16.1

### protocol buffers

The protocol buffer definitions are stored under `protobuf/*.proto`.

#### updating the schema

If editing the `*.proto` files, the Go bindings will need to be re-compiled. For example, the bindings for the sample message:

```
protoc --go_out=plugins=grpc:src/sample -I=protobuf protobuf/herald.proto
go generate
sh build-osx.sh
```

## Database

Sample records are stored via [bitcask db](https://pkg.go.dev/github.com/prologic/bitcask), which is currently hardcoded to live in `/tmp/db`.

Sample records are stored in a keyvalue store, where the label is the key and the sample protobuf message is the value.
