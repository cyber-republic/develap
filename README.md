## Pre-requisites
- Install Docker Engine
    - [Install Docker Engine for Ubuntu](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
    - [Install Docker Engine for Mac](https://docs.docker.com/docker-for-mac/install/)
    - [Install Docker Engine for Windows](https://docs.docker.com/docker-for-windows/install/)
    - For other platforms, check out [https://docs.docker.com/install/](https://docs.docker.com/install/)
- Golang v1.15 or higher
- Download the latest release specific to your platform from [https://github.com/cyber-republic/develap/releases](https://github.com/cyber-republic/develap/releases)

## URL Routing
- Mainchain Node RPC: 
```
http://localhost:5000/mainnet/mainchain
http://localhost:5000/testnet/mainchain
```
- DID Node RPC: 
```
http://localhost:5000/mainnet/did
http://localhost:5000/testnet/did
```
- ETH Node RPC: 
```
http://localhost:5000/mainnet/eth
http://localhost:5000/testnet/eth
```
- Cross-chain Transfer Service: 
```
http://localhost:5000/mainnet/xTransfer
http://localhost:5000/testnet/xTransfer
```
- Hive Node: 
```
http://localhost:5000/mainnet/hive
http://localhost:5000/testnet/hive
```

## How to Run
- Run a testnet environment with mainchain, did and eth nodes
    `./develap blockchain run -e testnet -n mainchain,did,eth`
- Run a testnet environment with mainchain node
    `./develap blockchain run -e testnet -n mainchain`
- Run a mainnet environment with did node
    `./develap blockchain run -e mainnet -n did`
- Kill eth node on mainnet environment
    `./develap blockchain kill -e mainnet -n eth`
- Kill all the nodes on testnet environment
    `./develap blockchain kill -e testnet`
- Show all the nodes currently running in mainnet environment
    `./develap blockchain list -e mainnet`

## Verify
If you run the mainchain node for mainnet, you can do the following to check whether it's working:
```
curl -H 'Content-Type: application/json' -H 'Accept:application/json' --data '{"method":"getcurrentheight"}' localhost:5000/mainnet/mainchain
```

## How to build the binary yourself
- `make` to build it for your local environment
- `make build-all` to build for 3 platforms: linux, darwin and windows

## Tools
- [build_dockerimages.sh](./tools/build_dockerimages.sh): This shell script automatically builds all the binaries for main chain, all the sidechains, services, etc and then packages them to be run inside docker images and if the flags "-p" and "-l" are set to "yes", the built docker images are automatically pushed to [Cyber Republic Docker Hub](https://cloud.docker.com/u/cyberrepublic/repository/list). Note that you need permission to push to the CR dockerhub but you can still build the images locally if you so choose
