## Pre-requisites
- Install Docker Engine
    - [Install Docker Engine for Ubuntu](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
    - [Install Docker Engine for Mac](https://docs.docker.com/docker-for-mac/install/)
    - [Install Docker Engine for Windows](https://docs.docker.com/docker-for-windows/install/)
    - For other platforms, check out [https://docs.docker.com/install/](https://docs.docker.com/install/)
- Golang v1.15 or higher
- Download the latest release specific to your platform from [https://github.com/cyber-republic/develap/releases](https://github.com/cyber-republic/develap/releases)

## URLs
- Status endpoints
``` 
http://localhost:5000/status/nodes
http://localhost:5000/status/nodes/running
http://localhost:5000/status/nodes/stopped
```
- Node endpoints: 
```
http://localhost:5000/node/mainnet/mainchain
http://localhost:5000/node/testnet/mainchain
http://localhost:5000/node/mainnet/did
http://localhost:5000/node/testnet/did
http://localhost:5000/node/mainnet/eth
http://localhost:5000/node/testnet/eth
```

## How to run different nodes
- Run a testnet environment with mainchain, did and eth nodes
    `./sws node run -e testnet -t mainchain,did,eth`
- Run a testnet environment with mainchain node
    `./sws node run -e testnet -t mainchain`
- Run a mainnet environment with did node
    `./sws node run -e mainnet -t did`
- Kill eth node on mainnet environment
    `./sws node kill -e mainnet -t eth`
- Kill all the nodes on testnet environment
    `./sws node kill -e testnet`
- Show all the nodes currently running in mainnet environment
    `./sws node list -e mainnet`

## How to aggregate all the running nodes into one unified URL
``` 
./sws serve
```

## Verify
If you run the did node for mainnet, you can do the following to check whether it's working:
```
curl -H 'Content-Type: application/json' -H 'Accept:application/json' --data '{"method":"getcurrentheight"}' localhost:5000/mainnet/did
```

## How to build the binary yourself
- `make` to build it for your local environment
- `make build-all` to build for 3 platforms: linux, darwin and windows

## Tools
- [build_dockerimages.sh](./tools/build_dockerimages.sh): This shell script automatically builds all the binaries for main chain, all the sidechains, services, etc and then packages them to be run inside docker images and if the flags "-p" and "-l" are set to "yes", the built docker images are automatically pushed to [Cyber Republic Docker Hub](https://cloud.docker.com/u/cyberrepublic/repository/list). Note that you need permission to push to the CR dockerhub but you can still build the images locally if you so choose
