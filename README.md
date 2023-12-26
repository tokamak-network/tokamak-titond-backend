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
- relayer : account for automatically claiming L2 -> L1 withdraws

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

## Test Deposit and Withdraw Guide