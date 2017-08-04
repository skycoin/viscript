## How to use

### Step 1

Add import to your application

`_ "github.com/skycoin/viscript/signal"`

This is all the code you need to change.

### Step 2

And then add arguments to the command line

```
run as a signal client(default:false)
-signal-client

connect to the address of signal server (default:"localhost:7999")
-signal-server-address

client id to use(default:1)
-signal-client-id
```

for example:

`~$ test localhost:20000 -signal-client -signal-client-id 2`