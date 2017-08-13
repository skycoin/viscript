// TODO: change readme
### Running with external app temporarily

1. Please make sure that the GOPATH is properly set and it's on the PATH.
2. Pull skycoin/skycoin repo
3. Using terminal:<br>
    ```cd skycoin/src/mesh/cmd/rpc/srv/```
4. Run following to install server in the `GOPATH/bin`:<br>
    ```go install```
5. Then navigate to the cli: <br>
    ```cd ../cli/```
6. Run following to install cli in the `GOPATH/bin`:<br>
    ```go install```
7. Go to the viscript main directory.
8. Run: <br>
    ```go install && viscript```

If everything's good you should be seeing two viscript terminals.
One that is in the back should be running meshnet server and one on top of that the meshnet cli.

NOTES:
- Please don't try to rename the binary files after go install of the server and cli.
    They are called "srv" and "cli" by default, for now.
- If you are wishing to run viscript rpc cli, please don't use `stp` command that starts new terminal
    with attached task.