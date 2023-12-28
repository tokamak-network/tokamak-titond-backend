# Titond Backend
![titond-architecture](./assets/titond%20architecture.png)
## Project overview

Titond is a PoC(Proof of concept) to provide on-demand L2. the goal of on-demand L2 is to make it easy for anyone to run their own Layer2 chain.

In this PoC, you can run titond for building L2 chain based on titan network.

This guide assumes that the user has knowledge of blockchain, kubernetes and aws.

you can also see korean version guide [here]()

**notice**
<br/>
**the current version does not support custom key. future update and you can run titond**

## Installation Instructions

### Prerequisites
**Golang**

Titond is developed in golang. you need to install golang version 1.20 or higher.
- [install golang](https://go.dev/doc/install)

**Ethereum API**

The titan network is L2 chain. L1 chain is required for L2 network to operate.

Titan network is based on the ethereum chain. we recommend using Sepolia test network. **if you use a ethereum mainnet, be aware that it uses a high amount of ether**

Select one of the node providers below and prepare sepolia network API access key

- [Infura](https://www.infura.io/)
- [Alchemy](https://www.alchemy.com/overviews/blockchain-node-providers)

**Kubernetes Cluster**

Titan network works on kubernetes. we use AWS EKS for the control plane and Fargete for the worker nodes.

but, AWS has charges and may not want to do this. you can build and use a local kubernetes cluster.

if you use a local kubernetes cluster, you need to modify manifests for your cluster environment.

- [Persistent Volume](./deployments/volume/pv.yaml)
- [L2geth Ingress](./deployments/l2geth/ingress.yaml)
- [Block Explorer Ingress](./deployments/explorer/ingress.yaml)

**Operator wallet**

To run titan network successfully, you need a few accounts. (To be update. for apply custom key)

- Deployer : deploy contracts that interacts with L1
- Sequencer : account for submits transaction batch to L1
- Proposer : account for submits proof to L1
- Block signer : account for operating titan network
- Relayer : account for automatically claiming L2 -> L1 withdraws

**AWS Resource**

Titond uses the AWS S3 service to store the L1 contract address and L2 genesis file.

You need to follow these steps:

1. Install AWS CLI and set AWS Configure
- [install AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
- [set AWS configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)

2. Create S3 bucket
- [create bucket](https://docs.aws.amazon.com/AmazonS3/latest/userguide/creating-bucket.html)

**Database for titond**

Titond uses PosgreSQL to store components information for the L2 chain.

- [Install PostgreSQL](https://www.postgresql.org/)

### Packages

| Titond          |                                         |
| --------------- | ----------------------------------------|
| http            | route request                           |
| api             | handle request                          |
| kubernetes      | control k8s cluster                     |            
|services         | used third party service (s3)           |
|db               | handle postgreSQL                       |
|model            | Database model                          |

### Commands

You can check the commands with the --help option.

```bash
$ go run cmd/titond/main.go --help

NAME:
   titond - The titond command line interface

USAGE:
   titond [global options] command [command options] [arguments...]

COMMANDS:
   check-swagger, g  check swagger
   help, h           Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --titond.contracts.rpc.url value         RPC url (default: "rpc-url") [$CONTRACTS_RPC_URL]
   --titond.contracts.deployer.key value    Deployer private key (default: "key") [$CONTRACTS_DEPLOYER_KEY]
   --titond.contracts.target.network value  Target network (default: "local") [$CONTRACTS_TARGET_NETWORK]
   --titond.sequencer.key value             SequencerKey (default: "titond") [$BATCH_SUBMITTER_SEQUENCER_PRIVATE_KEY]
   --titond.proposer.key value              ProposerKey (default: "titond") [$BATCH_SUBMITTER_PROPOSER_PRIVATE_KEY]
   --titond.relayer.wallet value            RelayerKey (default: "titond") [$MESSAGE_RELAYER__L1_WALLET]
   --titond.signer.key value                SignerKey (default: "titond") [$BLOCK_SIGNER_KEY]
   --kubernetes.incluster                   If this app is in cluster or not (default: false) [$KUBERNETES_INCLUSTER]
   --kubernetes.kubeconfig.path value       The path of kubeconfig (default: "$HOME/.kube/config") [$KUBERNETES_KUBECONFIG_PATH]
   --kubernetes.manifest.path value         The path of manifests (default: "./deployments") [$KUBERNETES_MANIFEST_PATH]
   --db.type value                          The type of database (default: "postgres") [$DB_TYPE]
   --db.host value                          The host of database (default: "localhost") [$DB_HOST]
   --db.port value                          The port of database (default: "5432") [$DB_PORT]
   --db.user value                          The user of database (default: "postgres") [$DB_USER]
   --db.password value                      The password of database (default: "postgres") [$DB_PASSWORD]
   --db.dbname value                        The database name of database (default: "titond") [$DB_DBNAME]
   --services.s3.bucket value               The S3 bucket name (default: "titond") [$SERVICES_S3_BUCKET]
   --services.s3.region value               The S3 region (default: "ap-northeast-2") [$SERVICES_S3_REGION]
   --http.host value                        The host of server (default: "0.0.0.0") [$HTTP_HOST]
   --http.port value                        The port of server (default: "8080") [$HTTP_PORT]
   --help, -h                               show help
```

## Usage Guide

You can input options and run titond

### Configuration

1. Check kubernetes config
    ```bash
    $ cat $HOME/.kube/config 
    ```
2. Set deployer key
    ```bash
    $ export CONTRACTS_DEPLOYER_KEY=<YOUR DEPLOYER PRIVATE KEY>
    ```
3. Set operator key
    ```bash
    $ export BATCH_SUBMITTER_SEQUENCER_PRIVATE_KEY=<YOUR SEQUENCER PRIVATE KEY>
    $ export BATCH_SUBMITTER_PROPOSER_PRIVATE_KEY=<YOUR PROPOSER PRIVATE KEY>
    $ export MESSAGE_RELAYER__L1_WALLET=<YOUR RELAYER PRIVATE KEY>
    $ export BLOCK_SIGNER_KEY=<YOUR BLOCK SIGNER PRIVATE KEY>
    ```
4. Set target network
    ```bash
    $ export CONTRACTS_TARGET_NETWORK=titond-demo 
    ```
5. Set L1 network RPC
    ```bash
    $ export CONTRACTS_RPC_URL=<YOUR L1 API ACCESS KEY>
    ```
6. Set database
    ```bash
    $ export DB_HOST=<YOUR DB HOST (default: localhost)>
    $ export DB_PORT=<YOUR DB PORT (default: 5432)>
    $ export DB_USER=<YOUR DB USER (default: postgres)>
    $ export DB_PASSWORD=<YOUR DB PASSWORD (default: postgress)>
    $ export DB_DBNAME=<YOUR DB NAME (default: titond)>
    ```
7. Set s3 bucket
    ```bash
    $ export SERVICES_S3_BUCKET=<YOUR BUCKET NAME (default: titond)>
    $ export SERVICES_S3_REGION=<YOUR BUCKET REGION (default: ap-northeast-2)>
    ```
8. Set host server
    ```bash
    $ export HTTP_HOST=<HOST SERVER IP (default: 0.0.0.0)
    $ export HTTP_PORT=<HOST SERVER PORT (default: 8080)>
    ```

### Run titond

Finished with the configuration, you can run titond
```bash
$ go run cmd/titond/main.go
```

### API endpoints

In this PoC version, we recommend only using APIs that POST and GET L2 chain components.

**Network**

The network is reponsible for deploying the contracts that will interact with the L2 on the L1 and creating the L2 genesis state file.

| Create network  |                                         |
| --------------- | ----------------------------------------|
| Description     | create network                          |
| Endpoint        | /api/networks                           |
| Http Method     | POST                                    |  

| Get network by id |                                         |
| ----------------  | ----------------------------------------|
| Description       | get network by ID                       |
| Endpoint          | /api/networks/{id}                      |
| Http Method       | GET                                     |

**Component**

The component is reponsible for building L2 chain.

| Create component|                                         |
| --------------- | ----------------------------------------|
| Description     | create component                        |
| Endpoint        | /api/components {body}                  |
| Http Method     | POST                                    | 

_Body_
```json
{
  "type": "component-type",
  "network_id": 1
}
```

type(string) : `data-transport-layer`, `l2geth`, `batch-submitter`, `relayer`, `explorer`
<br/>
network_id(int) : network id

| Get component by id |                                         |
| ----------------  | ------------------------------------------|
| Description       | get component by ID                       |
| Endpoint          | /api/components/{id}                      |
| Http Method       | GET                                       |

**To know detail, we provide API documentation as a swagger. you can access the documentation via the following URL.**

`<SERVER IP>:<SERVER PORT>/swagger/index.html`

_Tips of delete k8s objects_
<br/>
titond creates the namespace and deploys k8s objects on the namespace
<br/>
`namespace-{network id}`
<br/>
We have not provide DELETE APIs. but if you want to delete all of l2 chain components at once in your cluster, just delete the namespace
```bash
$ kubectl delete ns namespace-1
```

## Build L2 chain Guide

To build a Titan based L2 chain, you need the following essential components.

- l2geth: L2 client software.
- data-transport-layer: Indexing service for associated L1 data.
- batch-submitter: Service for sending L1 batched of transactions and results.

Titond also provides a few tools.

- explorer: L2 block explorer.
- relayer: a tool for automatically relaying L2 -> L1 messages.

To successfully build an L2 chain, you should follow these steps.

1. Deploy contracts to L1 and generate L2 genesis state file.
2. Create data-transport-layer.
3. Create l2geth.
4. Create batch-submitter.
5. Create relayer & explorer(selected)

### Build L2 chain

**1. Deploy contracts to L1 and generate L2 genesis state file**

**Create network**

**Action**

command:
```
$ curl -X POST {SERVER_IP}:{SERVER_PORT}/api/networks/
```

example:
```bash
$ curl -X POST localhost:8080/api/networks/
```

result:
```json
{
  "id": 1,
  "created_at": 1703729629,
  "updated_at": 1703729629,
  "contract_address_url": "",
  "state_dump_url": "",
  "status": false
}
```

It takes about 15-20 minutes for all contracts to be deployed on L1. when all of contracts are deployed, you can see the below:

**Get network**

**Action**

command:
```bash
$ curl -X GET {SERVER_IP}:{SERVER_PORT}/api/networks/{network_id}
```

example:
```bash
$ curl -X GET localhost:8080/api/networks/1
```

result:
```json
{
  "id": 1,
  "created_at": 1703729629,
  "updated_at": 1703729629,
  "contract_address_url": "{S3_BUCKET_URL}/address-1.json",
  "state_dump_url": "{S3_BUCKET_URL}/state-dump-1.json",
  "status": true
}
```

You can check the contents of that access the `contract_address_url` or `state_dump_url`.

contracts address:
```json
{
  "BondManager": "0x3a6f8FF68708194D7876CC7a71bb6033EA88200A",
  "Proxy__OVM_L1StandardBridge": "0xBA1AC8Bd0CdFBF619E899E7636525101c2dbb8A9",
  "Lib_AddressManager": "0xBEABF649E29bb050458686ca88B82B7e63A2381d",
  "CanonicalTransactionChain": "0xa0c08159d2c1492ae53a1e018735F3EF904970b8",
  "AddressDictator": "0x7C33d393f04945F5a6e2850D8Bf6Ad9FeBCCADaA",
  "StateCommitmentChain": "0x8bF6C4C5F8434a700DB159712979D58738F9AF31",
  "L1StandardBridge_for_verification_only": "0x9b0EAd143aB293d148B9E25b5c9A7F2dc035A32B",
  "Proxy__L1CrossDomainMessengerFast": "0xB36Cb435c91e969366A488De0814385E0317D8CC",
  "ChainStorageContainer-SCC-batches": "0x0F472f64C08337e1d3c062F5f4bbA06E8231Ad88",
  "Proxy__OVM_L1CrossDomainMessenger": "0xEdbc6709D9ecE42814B13cd97591261ACc549108",
  "ChugSplashDictator": "0xf72920892f672502a66641b8c00E3c704AE5A505",
  "OVM_L1CrossDomainMessenger": "0x3D47E7b02bA5Ae5Ea078ab41DbE21a39920B90BF",
  "ChainStorageContainer-CTC-batches": "0xf232619e16C0fFE86fBD06b34C7B406dC604FE9B",
  "L1CrossDomainMessengerFast": "0x227AC873632a81aeb509C23dA0077D111b8E4Bcf",
  "AddressManager": "0xBEABF649E29bb050458686ca88B82B7e63A2381d"
}
```

**2. Create data-transport-layer**

**Create data-transport-layer**

**Action**

command:
```bash
$ curl -d '{"type":"data-transport-layer", "network_id":{nerwork_id}}' \ 
-H "Content-Type: application/json" \
-X POST {SERVER_IP}:{SERVER_PORT}/api/components/
```

example:
```bash
$ curl -d '{"type":"data-transport-layer", "network_id":1 }' \
-H "Content-Type: application/json" \
-X POST localhost:8080/api/components/
```

result:
```json
{
  "id": 1,
  "created_at": 1703846585,
  "updated_at": 1703846585,
  "name": "dtl-pod-name",
  "type": "data-transport-layer",
  "status": false,
  "public_url": ""
}
```

It takes about 10-15 minutes to be created.

**Get data-transport-layer**

After data-transport-layer is finished successfully, you can check below result.

**Action**

command:
```bash
$ curl -X GET {SERVER_IP}:{SERVER_PORT}/api/components/{component_id}
```

example:
```bash
$ curl -X GET localhost:8080/api/components/1
```

result:
```json
{
  "id": 1,
  "created_at": 1703846585,
  "updated_at": 1703846585,
  "name": "dtl-pod-name",
  "type": "data-transport-layer",
  "status": true,
  "public_url": ""
}
```

**3. Create l2geth**

**Create l2geth**

**Action**

command:
```bash
$ curl -d '{"type":"l2geth", "network_id":{nerwork_id}}' \ 
-H "Content-Type: application/json" \
-X POST {SERVER_IP}:{SERVER_PORT}/api/components/
```

example:
```bash
$ curl -d '{"type":"l2geth", "network_id":1 }' \
-H "Content-Type: application/json" \
-X POST localhost:8080/api/components/
```

result:
```json
{
  "id": 2,
  "created_at": 1703894525,
  "updated_at": 1703894525,
  "name": "l2geth-pod-name",
  "type": "l2geth",
  "status": false,
  "public_url": ""
}
```

It takes about 1-5 minutes to be created.

**Get l2geth**

After l2geth is finished successfully, you can check below result.

**Action**

command:
```bash
$ curl -X GET {SERVER_IP}:{SERVER_PORT}/api/components/{component_id}
```

example:
```bash
$ curl -X GET localhost:8080/api/components/2
```

result:
```json
{
  "id": 2,
  "created_at": 1703894525,
  "updated_at": 1703894525,
  "name": "l2geth-pod-name",
  "type": "l2geth",
  "status": true,
  "public_url": "l2geth-1-2.titond-holesky.tokamak.network"
}
```

**4. Create batch-submitter**

**Create batch-submitter**

**Action**

command:
```bash
$ curl -d '{"type":"batch-submitter", "network_id":{nerwork_id}}' \ 
-H "Content-Type: application/json" \
-X POST {SERVER_IP}:{SERVER_PORT}/api/components/
```

example:
```bash
$ curl -d '{"type":"batch-submitter", "network_id":1 }' \
-H "Content-Type: application/json" \
-X POST localhost:8080/api/components/
```

result:
```json
{
  "id": 3,
  "created_at": 1703984245,
  "updated_at": 1703984245,
  "name": "batch-submitter-pod-name",
  "type": "batch-submitter",
  "status": false,
  "public_url": ""
}
```

It takes about 1-5 minutes to be created.

**Get batch-submitter**

After batch-submitter is finished successfully, you can check below result.

**Action**

command:
```bash
$ curl -X GET {SERVER_IP}:{SERVER_PORT}/api/components/{component_id}
```

example:
```bash
$ curl -X GET localhost:8080/api/components/1
```

result:
```json
{
  "id": 3,
  "created_at": 1703984245,
  "updated_at": 1703984245,
  "name": "batch-submitter-pod-name",
  "type": "batch-submitter",
  "status": true,
  "public_url": ""
}
```

**5. Create relayer & explorer(selected)**

**Create relayer & explorer**

You can create a relayer and an explorer parallel

**Action**

command:
```bash
//relayer
$ curl -d '{"type":"relayer", "network_id":{nerwork_id}}' \ 
-H "Content-Type: application/json" \
-X POST {SERVER_IP}:{SERVER_PORT}/api/components/

//explorer
$ curl -d '{"type":"explorer", "network_id":{nerwork_id}}' \ 
-H "Content-Type: application/json" \
-X POST {SERVER_IP}:{SERVER_PORT}/api/components/
```

example:
```bash
//relayer
$ curl -d '{"type":"relayer", "network_id":1 }' \
-H "Content-Type: application/json" \
-X POST localhost:8080/api/components/

//explorer
$ curl -d '{"type":"explorer", "network_id":1 }' \
-H "Content-Type: application/json" \
-X POST localhost:8080/api/components/
```

result:
```json
//relayer
{
  "id": 5,
  "created_at": 1704092938,
  "updated_at": 1704092938,
  "name": "relayer-pod-name",
  "type": "relayer",
  "status": false,
  "public_url": ""
}

//explorer
{
  "id": 6,
  "created_at": 1704104523,
  "updated_at": 1704104523,
  "name": "explorer-pod-name",
  "type": "explorer",
  "status": false,
  "public_url": ""
}
```

It takes about 2-8 minutes to be created.

**Get relayer & explorer**

After relayer & explore is finished successfully, you can check below result.

**Action**

command:
```bash
$ curl -X GET {SERVER_IP}:{SERVER_PORT}/api/components/{component_id}
```

example:
```bash
$ curl -X GET localhost:8080/api/components/5
```

result:
```json
//relayer
{
    "id": 5,
    "created_at": 1704092938,
    "updated_at": 1704092938,
    "name": "relayer-pod-name",
    "type": "relayer",
    "status": true,
    "public_url": ""
}

//explorer
{
    "id": 6,
    "created_at": 1704104523,
    "updated_at": 1704104523,
    "name": "explorer-pod-name",
    "type": "explorer",
    "status": true,
    "public_url": "explorer-1-6.titond-holesky.tokamak.network"
}
```

## Deposit and Withdraw Guide

### Connect L2 chain your crypto wallets

You can connect the L2 chain you built in a crypto wallet like metamask

| New L2 network  |                                                                               |
| --------------- | ------------------------------------------------------------------------------|
| RPC URL         | your l2geth public url(ex. https://l2geth-1-2.titond-holesky.tokamak.network) |
| Chain ID        | 17                                                                            |
| Currency symbol | ETH                                                                           |            

### Deposit

if you build l2chain successfully, you can deposit ETH from L2 to L1.

1. Check L1 contracts address in `{S3_BUCKET_URL}/address-1.json`
2. You can find `Proxy__OVM_L1StandardBridge` address
    ```json
    {
      ...

      "Proxy__OVM_L1StandardBridge": "0xBA1AC8Bd0CdFBF619E899E7636525101c2dbb8A9",

      ...
    }
    ```
3. Change network to L1, send some ether to `Proxy__OVM_L1StandardBridge`

### Withdraw

TBD