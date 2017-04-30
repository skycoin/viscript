#!/bin/bash

# Set verbosity variable to false
let v=false

# Check if -v is passed and if true set verbosity to true to print everything
if [ "$1" = "-v" ]; then
    v=true
fi

# Define print verbose function that looks at the verbosity option and echoes
pv () {
    if [ $v = true ]; then
        echo "[ • ]" $1
    fi
}

echo "Start of packaging viscript and binaries:"

# Set Root directory name
readonly ROOT_DIR="$PWD/Viscript"
pv "Setting root directory to $ROOT_DIR"

# Make directory for Viscript
if [ ! -d $ROOT_DIR ]; then
    mkdir $ROOT_DIR
    pv "Creating root directory"
else
    pv "Root directory already exists, cleaning it up."
    rm -rf $ROOT_DIR 
    mkdir $ROOT_DIR
fi

# Make bin folder inside of ROOT_DIR
pv "Creating bin in root directory"
mkdir $ROOT_DIR/bin

# Make skycoin folder inside of ROOT_DIR/bin
pv "Creating skycoin directory inside root/bin"
mkdir $ROOT_DIR/bin/skycoin

# Make meshnet folder inside of ROOT_DIR/bin
pv "Creating meshnet dir inside root/bin"
mkdir $ROOT_DIR/bin/meshnet

# Set the Skycoin path from github
githubSkycoinPath="github.com/skycoin/skycoin"

# Go get skycoin
pv "Go getting Skycoin: $githubSkycoinPath"
go get -u -d $githubSkycoinPath &>/dev/null

# Set local skycoin path
localSkycoinPath="$GOPATH/src/$githubSkycoinPath"
pv "Local Skycoin path set to: $localSkycoinPath"

# Check if skycoin directory exists in gopath
if [ ! -d "$localSkycoinPath" ]; then
    pv "Skycoin directory doesn't exist in $GOPATH/src/github.com/"
    exit 1
else
    pv "Skycoin directory verified" 
fi

# Change Directory to local skycoin path
pv "Changing directory to Skycoin main file"
cd "$localSkycoinPath/cmd/skycoin/"

# Get all dependencies for skycoin.go
pv "Getting all dependencies for Skycoin"
go get -d ./...

# Build skycoin.go 
pv "Building Skycoin binary"
go build -o skycoin skycoin.go

# Check if skycoin binary was built successfully
if [ ! -f "skycoin" ]; then
    pv "Building Skycoin binary failed. Exiting"
    exit 1
fi

# Create skycoin directory inside the root/bin 
pv "Copying skycoin binary inside $ROOT_DIR/bin/skycoin/"
mv skycoin $ROOT_DIR/bin/skycoin/

# Change to src/gui/static 
pv "Changing directory to src/gui"
cd "$localSkycoinPath/src/gui/"

# Make directory for skycoin statics
pv "Making directory static inside local skycoin/"
mkdir "$ROOT_DIR/bin/skycoin/static/"

# Copy static folder to the local skycoin path inside newly created static/
pv "Copying static for Skycoin"
cp -R "static/" "$ROOT_DIR/bin/skycoin/"

# Change directory to local Skycoin cli 
pv "Changing directory to Skycoin cli"
cd "$localSkycoinPath/cmd/cli/"

# Get all dependencies for cli.go
pv "Getting all dependencies for Skycoin cli"
go get -d ./...

# Build cli.go
pv "Building Skycoin cli"
go build -o skycoin-cli cli.go

# Check if building cli was successfull
if [ ! -f "skycoin-cli" ]; then
    pv "Building Skycoin cli failed. Exiting"
    exit 1
fi

# Move Skycoin cli to skycoin path
pv "Moving Skycoin cli to the root directory"
mv skycoin-cli "$ROOT_DIR/bin/skycoin/"

# Change directory to local Skycoin mesh server
pv "Changing directory to Skycoin rpc server"
cd "$localSkycoinPath/src/mesh/cmd/rpc/srv/"

# Get all dependencies for rpc_run.go
pv "Getting all dependencies for meshnet server"
go get -d ./...

# Build run_rpc.go
pv "Building Skycoin rpc server"
go build -o meshnet-server rpc_run.go

# Check if building cli was successfull
if [ ! -f "meshnet-server" ]; then
    pv "Building Skycoin rpc server failed. Exiting"
    exit 1
fi

# Move Skycoin mesh server to root dir
pv "Moving Skycoin rpc server to the root directory"
mv meshnet-server "$ROOT_DIR/bin/meshnet/"

# Change directory to local Skycoin mesh cli
pv "Changing directory to Skycoin rpc mesh cli"
cd "$localSkycoinPath/src/mesh/cmd/rpc/cli/"

# Get all dependencies for meshnet cli
pv "Getting all dependencies for meshnet cli"
go get -d ./...

# Build cli.go
pv "Building Skycoin rpc mesh cli"
go build -o meshnet-cli cli.go

# Check if building cli was successfull
if [ ! -f "meshnet-cli" ]; then
    pv "Building Skycoin rpc mesh cli failed. Exiting"
    exit 1
fi

# Move Skycoin mesh cli to root dir
pv "Moving Skycoin rpc mesh cli to the root directory"
mv meshnet-cli "$ROOT_DIR/bin/meshnet/"

# Change directory to root
pv "Changing directory to current working one"
cd "$ROOT_DIR/" && cd ..

# Build viscript.go
pv "Building viscript binary"
go build -o viscript viscript.go

# Check if building viscript binary was successfull
if [ ! -f "viscript" ]; then
    pv "Building viscript binary failed. Exiting"
    exit 1
fi

# Move viscript binary to the root
pv "Moving viscript binary into the root directory"
mv "viscript" "$ROOT_DIR/"

# Copy README.md to the root
pv "Copying README.md for viscript"
cp "README.md" "$ROOT_DIR/"

# Move all assets that are required for viscript
pv "Copying assets folder inside the root directory for viscript"
cp -R "assets/" "$ROOT_DIR/"

# Build viscript cli
pv "Building viscript cli"
go build -o viscript-cli rpc/cli/cli.go

# Check if building cli was successfull
if [ ! -f "viscript-cli" ]; then
    pv "Building viscript cli failed. Exiting"
    exit 1
fi

# Move viscript cli to the root directory 
pv "Moving viscript cli to the root directory"
mv viscript-cli "$ROOT_DIR/"

# Create viscript-cli.sh file that uses gotty to run the cli
pv "Creating bash file to run cli with gotty: https://github.com/yudai/gotty"
gottyCommand="gotty -w -p 9999 --reconnect ./viscript-cli"
echo $gottyCommand > "$ROOT_DIR/viscript-cli.sh" 

# TODO: zip here?
