> note: this is incomplete documentation and is more of a brain dump whilst herald is being developed...

## Processes

**Herald** tags samples and experiments with _processes_. Processes included "sequencing", "basecalling", "pipeline X" etc.

### Adding a process

> note: the UI will only be updated if the process is for samples (I've left the experiment processes harcoded in the UI for now)

1. add a process definition in `processes.go`

This file is located in the source code (`src/data/processes.go`). An example function call is:

```go
createProcessDefinition("mypipeline", []string{"sequence", "basecall"}, false, true)
```

This will create a process called `mypipeline` which depends on the processes `sequence` and `basecall` being complete.

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
protoc -I=protobuf --go_out=plugins=grpc:src/ protobuf/*.proto
go generate
sh build-osx.sh
```

To update the MinKNOW api:

```
protoc --go_out=plugins=grpc:src/ont_rpc -I protobuf/minknow/rpc protobuf/minknow/rpc/rpc_options.proto
```

## Database

Sample records are stored via [bitcask db](https://pkg.go.dev/git.mills.io/prologic/bitcask), which is currently hardcoded to live in `/tmp/db`.

Sample records are stored in a keyvalue store, where the label is the key and the sample protobuf message is the value.
