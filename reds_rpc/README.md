## Running CLI

In order to be able to run CLI client for viscript RPC server should be running.
At this point simple RPC server is run using goroutine with the viscript runtime when installing
and running the viscript:

```
go install && viscript
```

Once the server is listening we can just type:

```
go run rpc/cli/cli.go
```

And the CLI connects to the RPC server which is ready to serve.

## Running with gotty

You could also try launching `cli.sh` script like:

```
chmod +x cli.sh
./cli.sh
```
to be able to interact with cli using the web interface (uses gotty).