# Titond Backend
![titond-architecture](./assets/titond%20architecture.png)
## Project overview

타이톤드는 on-demand L2를 제공하기 위한 PoC 입니다. on-demand L2 는 자신만의 Layer2 체인을 누구든지 쉽게 만들수 있는 것을 목표로 합니다.

이번 PoC 에서는 titond로 타이탄 네트워크 기반의 L2 체인을 구축할 수 있습니다.

이 가이드는 블록체인, 쿠버네티스, aws에 대한 기본적인 지식이 있는 유저를 상대로 진행됩니다.

**notice**
<br/>
**현재 버전은 개인키 커스텀을 지원하지 않습니다. 향후 업데이트 되면 titond를 실행할 수 있습니다**

## Installation Instructions

### Prerequisites
**Golang**

타이톤드는 고 언어로 개발됐습니다. 1.20 혹은 그 이상의 고 언어 버전이 필요합니다.
- [install golang](https://go.dev/doc/install)

**Ethereum API**

타이탄 네트워크는 L2 체인입니다. L2 네트워크를 운영하기 위해서는 L1 체인이 필요합니다.

타이탄 네트워크는 이더리움을 기반으로 하고 있습니다. 세폴리아(Sepolia) 테스트 네트워크를 사용하는 것을 권장합니다. **만약 이더리움 메인넷을 사용할 경우 많은 양의 이더가 사용될 수 있으니 주의하세요**

아래 노드 프로바이더 중 하나를 선택해서 세폴리아 네트워크 API 키를 준비하세요.

- [Infura](https://www.infura.io/)
- [Alchemy](https://www.alchemy.com/overviews/blockchain-node-providers)

**Kubernetes Cluster**

타이탄 네트워크는 쿠버네티스 위에서 동작합니다. 현재 우리는 컨트롤 플레인으로 AWS EKS를 사용하고 워커 노드로 Fargate를 사용합니다.

그러나, AWS는 과금이 되고 이를 원하지 않을수 있습니다. 그래서 로컬 쿠버네티스를 구축하고 사용할 수 있습니다.

만약 로컬 쿠버네티스 클러스터를 구축해서 사용한다면 몇가지 manifest 파일을 로컬 쿠버네티스 클러스터 환경에 맞게 수정이 필요합니다.

- [Persistent Volume](./deployments/volume/pv.yaml)
- [L2geth Ingress](./deployments/l2geth/ingress.yaml)
- [Block Explorer Ingress](./deployments/explorer/ingress.yaml)

**Operator wallet**

타이탄 네트워크를 성공적으로 구축하기 위해서는 몇가지 계정(EOA)이 필요합니다. (개인키 커스텀을 지원하기 위해 업데이트 될 예정.)

- Deployer : L1 네트워크와 상호작용 하기 위한 컨트랙트를 배포합니다.
- Sequencer : L1에 트랜잭션 배치를 제출하기 위한 계정
- Proposer : L1에 proof를 제출하기 위한 계정
- Block signer : 타이탄 네트워크를 운영하기 위한 계정
- Relayer : L2 -> L1 으로 출금을 자동으로 클레임 하기 위한 계정

**AWS Resource**

타이톤드는 L1 컨트랙트 주소와 L2 제네시스 파일을 저장하기 위해 AWS S3를 사용합니다.

아래 단계를 따라주세요:

1. Install AWS CLI and set AWS Configure
- [install AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
- [set AWS configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)

2. Create S3 bucket
- [create bucket](https://docs.aws.amazon.com/AmazonS3/latest/userguide/creating-bucket.html)

**Database for titond**

타이톤드는 L2 체인 컴포넌트 정보를 저장하기 위해 PostgreSQL 를 사용합니다.

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

--help 옵션을 통해 명령어를 확인할 수 있습니다.

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

옵션을 입력하고 타이톤드를 실행할 수 있습니다.

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

설정이 모두 완료되면 타이톤드를 실행할 수 있습니다.
```bash
$ go run cmd/titond/main.go
```

### API endpoints

PoC 버전에서는 POST와 GET API만 사용할 것을 권장합니다.

**Network**

네트워크는 L2와 상호작용하기 위한 컨트랙트를 L1에 배포하고 L2 제네시스 파일을 만듭니다.

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

컴포넌트는 L2체인을 구축하기 위한 구성요소 입니다.

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

**타이톤드는 swagger로 API 문서를 제공합니다. 자세하게 알고싶다면 아래 URL을 통해서 문서에 접근할 수 있습니다.**

`<SERVER IP>:<SERVER PORT>/swagger/index.html`

_쿠버네티스 오브젝트 삭제 팁_
<br/>
타이톤드는 쿠버네티스 오브젝트를 배포하기 위해 네임스페이스를 만듭니다.
<br/>
`namespace-{network id}`
<br/>
현재 DELETE API를 제공하지 않습니다. 그러나 클러스터에서 l2 체인 컴포넌트들을 한번에 모두 지우고 싶다면 네임스페이스만 지우면 됩니다.
```bash
$ kubectl delete ns namespace-1
```

## Build L2 chain Guide

L2 체인인 타이탄을 구축하기 위해서는 다음과 같은 필수적인 컴포넌트들이 필요합니다.

- l2geth: L2 클라이언트 소프트웨어.
- data-transport-layer: 연결된 L1 데이터에 대한 인덱싱 서비스.
- batch-submitter: 트랜잭션 배치를 L1으로 전송해주기 위한 서비스.

타이톤드는 몇가지 툴도 제공합니다.

- explorer: L2 블록 탐색기.
- relayer: L2 -> L1 메세지를 자동으로 전달해주기 위한 툴.

L2 체인을 성공적으로 구축하기 위해서는 다음과 같은 단계를 따라야 합니다.

1. L1에 컨트랙트를 배포하고 L2 제네시스 파일 생성.
2. data-transport-layer 생성.
3. l2geth 생성.
4. batch-submitter 생성.
5. relayer & explorer(선택사항) 생성.

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

L1에 컨트랙트를 모두 배포하기 위해서는 대략 15-20분 정도 소요됩니다. 컨트랙트가 모두 배포되면 아래와 같은 결과를 볼 수 있습니다:

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

`contract_address_url` or `state_dump_url` 에 접근해서 내용을 확인할 수 있습니다.

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

완료까지 10-15분 정도 소요됩니다.

**Get data-transport-layer**

data-transport-layer가 성공적으로 만들어진 뒤 아래와 같은 결과를 확인할 수 있습니다.

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

완료까지 1-5분 정도 소요됩니다.

**Get l2geth**

l2geth가 성공적으로 만들어진 뒤 아래와 같은 결과를 확인할 수 있습니다.

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

완료까지 1-5분 정도 소요됩니다.

**Get batch-submitter**

bath-submitter가 성공적으로 만들어진 뒤 아래와 같은 결과를 확인할 수 있습니다.

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

relayer와 explorer는 동시에 생성이 가능합니다.

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

완료까지 2-8분 정도 소요됩니다.

**Get relayer & explorer**

relayer와 explore가 성공적으로 만들어진 뒤 아래와 같은 결과를 확인할 수 있습니다.

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

메타마스크와 같은 암호화폐 지갑에서 구축한 L2 체인을 연결할 수 있습니다.

| New L2 network  |                                                                               |
| --------------- | ------------------------------------------------------------------------------|
| RPC URL         | your l2geth public url(ex. https://l2geth-1-2.titond-holesky.tokamak.network) |
| Chain ID        | 17                                                                            |
| Currency symbol | ETH                                                                           |            

### Deposit

L2 체인을 성공적으로 구축했다면 이더를 L2에서 L1으로 입금할 수 있습니다.

1. `{S3_BUCKET_URL}/address-1.json` 에서 L1에 배포된 컨트랙트 주소 확인
2. `Proxy__OVM_L1StandardBridge` 주소 확인
    ```json
    {
      ...

      "Proxy__OVM_L1StandardBridge": "0xBA1AC8Bd0CdFBF619E899E7636525101c2dbb8A9",

      ...
    }
    ```
3. 지갑에 연결된 네트워크를 L1 네트워크로 바꾼 뒤 `Proxy__OVM_L1StandardBridge`으로 이더를 전송

### Withdraw

TBD