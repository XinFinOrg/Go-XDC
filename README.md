# go-xdc
A wrapper CLI for setting up the xdc blockchain network which uses the [XDC01 Docker Nnodes](https://github.com/XinFinorg/XDC01-docker-Nnodes)

Current Version :- Quorum v2.1.0

## Prerequisite
**Operating System**: Ubuntu 16.04 64-bit or higher

**Tools**: Go v1.10 or higher with approriate PATH variables set
 
## Clone repository in your $GOPATH

    cd $GOPATH/src/github.com/yourusername/

    git clone https://github.com/XinFinOrg/Go-XDC.git
    
    cd Go-XDC

## Install Dependencies 
    go get ./...

## Build
    go build -o xdc
    

## Run
    ./xdc

```
Usage:
  xdc [command]

Available Commands:
  help        Help about any command
  prepare     Installs the dependencies for network setup
  setup       Setup XDC network
  start       Start the XDC network
  stop        Stop the XDC network

Flags:
      --config string   config file (default is $HOME/.go-xdc.yaml)
  -h, --help            help for xdc
  -t, --toggle          Help message for toggle

Use "xdc [command] --help" for more information about a command.

Complete documentation available at https://www.xinfin.org/

```
