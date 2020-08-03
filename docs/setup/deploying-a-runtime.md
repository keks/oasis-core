# Deploying a Runtime

Before proceeding, make sure to look at the [prerequisites] required for running
an Oasis Core environment followed by [build instructions] for the respective
environment (non-SGX or SGX), using the [oasis net runner] and see
[runtime documentation] for a general documentation of runtimes.

These instructions will show how to register and deploy a runtime node on a
local development network.

[prerequisites]: prerequisites.md
[build instructions]: building.md
[oasis net runner]: oasis-net-runner.md
[runtime documentation]: ../runtime/index.md

## Provision a Single Validator Node Network

Use the [oasis net runner] to provision a validator node network without any
registered runtime.

<!-- markdownlint-disable line-length -->
```
mkdir /tmp/runtime-example

oasis-net-runner \
  --basedir.no_temp_dir \
  --basedir /tmp/runtime-example \
  --fixture.default.node.binary go/oasis-node/oasis-node \
  --fixture.default.setup_runtimes=false \
  --fixture.default.deterministic_entities \
  --fixture.default.num_entities 2 \
  --fixture.default.fund oasis1qznshq4ttrgh83d9wqvgmsuq3pfsndg3tus7lx98
```
<!-- markdownlint-enable line-length -->

The following steps should be run in a separate terminal window.
To simplify the instructions set up an `ADDR` environment variable
pointing to the UNIX socket exposed by the started node:

```
export ADDR=unix:/tmp/runtime-example/net-runner/network/validator-0/internal.sock
```

Confirm the network is running by listing all registered entities:

```
oasis-node registry entity list -a $ADDR -v
```

Should give output similar to:

<!-- markdownlint-disable line-length -->
```
{"v":1,"id":"JTUtHd4XYQjh//e6eYU7Pa/XMFG88WE+jixvceIfWrk=","nodes":["LQu4ZtFg8OJ0MC4M4QMeUR7Is6Xt4A/CW+PK/7TPiH0="]}
{"v":1,"id":"+MJpnSTzc11dNI5emMa+asCJH5cxBiBCcpbYE4XBdso="}
{"v":1,"id":"TqUyj5Q+9vZtqu10yw6Zw7HEX3Ywe0JQA9vHyzY47TU=","allow_entity_signed_nodes":true}
```
<!-- markdownlint-enable line-length -->

In following steps we will register and run the [simple-keyvalue] runtime on the
network.

[simple-keyvalue]: ../../tests/runtimes/simple-keyvalue

## Initializing a Runtime

To generate and sign a runtime registration transaction that will initialize and
register the runtime we will use the `registry runtime gen_register` command.
When initializing a runtime we need to specify various runtime parameters. To
list all of the available parameters with short descriptions run:
`registry runtime gen_register --help` subcommand.

```
oasis-node registry runtime gen_register --help
```

For additional information about runtimes and parameters see the
[runtime documentation] and [code reference].

Before generating the registration transaction, gather the following data and
set up environment variables to simplify instructions.

- `ENTITY_DIR` - Path to the entity directory created when starting the
development network. This entity will be the runtime owner. The genesis used in
the provisioning initial network step funds the `entity-2` entity (located
within the net-runner directory), so that is the one we use in following
instructions.
- `ENTITY_ID` - ID of the entity that will be the owner of the runtime. You can
get the entity ID from `$ENTITY_DIR/entity.json` file.
- `GENESIS_JSON` - Path to the genesis.json file used in the development
network.
- `RUNTIME_ID` - See [runtime identifiers] on how to choose a runtime
identifier. In the example we use
 `8000000000000000000000000000000000000000000000000000000001234567`.
- `NONCE` - Entity account nonce. If you followed the guide, nonce `0`
would be the initial nonce to use for the entity. Note: make sure to keep
updating the nonce when generating new transactions. To query for current
account nonce value use [stake account info] CLI.

```
export ENTITY_DIR=/tmp/runtime-example/net-runner/network/entity-2/
export ENTITY_ID=+MJpnSTzc11dNI5emMa+asCJH5cxBiBCcpbYE4XBdso=
export GENESIS_JSON=/tmp/runtime-example/net-runner/network/genesis.json
export RUNTIME_ID=8000000000000000000000000000000000000000000000000000000001234567
export NONCE=0
```

[runtime identifiers]: ../runtime/identifiers.md
[stake account info]:  ../oasis-node/cli.md#info

```
oasis-node registry runtime gen_register \
  --transaction.fee.gas 1000 \
  --transaction.fee.amount 1000 \
  --transaction.file /tmp/runtime-example/register_runtime.tx \
  --transaction.nonce $NONCE \
  --genesis.file $GENESIS_JSON \
  --signer.backend file \
  --signer.dir $ENTITY_DIR \
  --runtime.id $RUNTIME_ID \
  --runtime.kind compute \
  --runtime.executor.group_size 1 \
  --runtime.merge.group_size 1 \
  --runtime.txn_scheduler.group_size 1 \
  --runtime.storage.group_size 1 \
  --runtime.admission_policy entity-whitelist \
  --runtime.admission_policy_entity_whitelist $ENTITY_ID \
  --debug.dont_blame_oasis \
  --debug.allow_test_keys
```

This command outputs a signed transaction in the
`/tmp/runtime-example/register_runtime.tx` file. In the next step we will submit
the transaction to complete the runtime registration.

NOTE: when registering a runtime on a non-development network you will likely
want to modify default parameters. Additionally, since we are running this on
a test network, we need to enable the debug flags above.

<!-- markdownlint-disable line-length -->
[code reference]: https://pkg.go.dev/github.com/oasisprotocol/oasis-core/go/registry/api?tab=doc#Runtime
<!-- markdownlint-enable line-length -->

## Submitting the Runtime Register Transaction

To register the runtime, submit the generated transaction.

```
oasis-node consensus submit_tx \
    --transaction.file /tmp/runtime-example/register_runtime.tx \
    --address $ADDR
```

## Confirm Runtime is Registered

To confirm the runtime is registered use the
`registry runtime list -v --include_suspended` command.

```
oasis-node registry runtime list -a $ADDR -v --include_suspended
```

Should give output similar to

<!-- markdownlint-disable line-length -->
```
{"v":1,"id":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEjRWc=","entity_id":"JTUtHd4XYQjh//e6eYU7Pa/XMFG88WE+jixvceIfWrk=","genesis":{"state_root":"xnK40e9W7Sirh8NiLFEUBpvdOte4+XN0mNDAHs7wlno=","state":null,"storage_receipts":null,"round":0},"kind":1,"tee_hardware":0,"versions":{"version":{"Major":0,"Minor":0,"Patch":0}},"executor":{"group_size":1,"group_backup_size":0,"allowed_stragglers":0,"round_timeout":10000000000},"merge":{"group_size":1,"group_backup_size":0,"allowed_stragglers":0,"round_timeout":10000000000},"txn_scheduler":{"group_size":1,"algorithm":"batching","batch_flush_timeout":1000000000,"max_batch_size":1000,"max_batch_size_bytes":16777216},"storage":{"group_size":1,"min_write_replication":1,"max_apply_write_log_entries":100000,"max_apply_ops":2,"max_merge_roots":8,"max_merge_ops":2,"checkpoint_interval":0,"checkpoint_num_kept":0,"checkpoint_chunk_size":0},"admission_policy":{"any_node":{}},"staking":{}}
```
<!-- markdownlint-enable line-length -->

Note: since we did not setup any runtime nodes, the runtime will get
[suspended] until nodes for the runtime register.

In the next step we will setup and run a runtime node.

[suspended]: ../runtime/index.md#Suspending-Runtimes

## Running a Runtime Node

We will now run a node that will act as a compute, storage and client node for
the runtime.

NOTE: in a real deployment scenario each of the compute, storage and client
types of nodes would likely be a separate deployment of nodes.

Before running the node gather the following data parameters and set up
environment variables to simplify instructions.

- `RUNTIME_BINARY` - Path to the runtime binary that will be run on the node.
We will use the [simple-keyvalue] runtime. If you followed the
[build instructions] the built binary is available at
./target/default/debug/simple-keyvalue.
- `SEED_NODE_ADDRESS` - Address of the seed node in the development network.

<!-- markdownlint-disable line-length -->
```
export RUNTIME_BINARY=/workdir/target/default/debug/simple-keyvalue
export SEED_NODE_ADDRESS=2F976E6C38D0DD73C046BDA1A6555EBF65BE6CC9@127.0.0.1:20002

mkdir -m 0700 /tmp/runtime-example/runtime-node

oasis-node \
  --datadir /tmp/runtime-example/runtime-node \
  --log.level debug \
  --log.format json \
  --log.file /tmp/runtime-example/runtime-node/node.log \
  --worker.registration.entity $ENTITY_DIR/entity.json \
  --genesis.file $GENESIS_JSON \
  --storage.backend badger \
  --worker.storage.enabled \
  --worker.compute.enabled \
  --worker.runtime.provisioner unconfined \
  --worker.txn_scheduler.check_tx.enabled \
  --grpc.log.debug \
  --runtime.supported $RUNTIME_ID \
  --runtime.history.tag_indexer.backend bleve \
  --worker.runtime.paths $RUNTIME_ID=$RUNTIME_BINARY \
  --tendermint.debug.addr_book_lenient \
  --tendermint.debug.allow_duplicate_ip \
  --debug.dont_blame_oasis \
  --debug.allow_test_keys \
  --tendermint.p2p.seed $SEED_NODE_ADDRESS
```
<!-- markdownlint-enable line-length -->

**NOTE: This also enables unsafe debug-only flags which must never be used in a
production setting as they may result in node compromise.**

Following steps should be run in a new terminal window.

## Updating Entity Nodes

Before the newly started runtime node can register itself as a runtime node, we
need to update the entity information, to include the started node.

Before proceeding gather the runtime node id and store it in a variable. If you
followed above instructions, the node id can be seen in
`/tmp/runtime-example/runtime-node/identity_pub.pem`.

Update the entity and generate a transaction that will update the registry
state.

```
export NODE_ID=vWUfSmjrHSlN5tSSO3/Qynzx+R/UlwPV9u+lnodQ00c=

oasis-node registry entity update \
  --signer.dir $ENTITY_DIR  \
  --entity.node.id $NODE_ID

oasis-node registry entity gen_register \
  --genesis.file $GENESIS_JSON \
  --signer.backend file \
  --signer.dir $ENTITY_DIR \
  --transaction.file /tmp/runtime-example/update_entity.tx \
  --transaction.fee.gas 2000 \
  --transaction.fee.amount 2000 \
  --transaction.nonce $NONCE \
  --debug.dont_blame_oasis \
  --debug.allow_test_keys
```

Submit the generated transaction:

```
oasis-node consensus submit_tx \
    --transaction.file /tmp/runtime-example/update_entity.tx \
    --address $ADDR
```

Confirm the entity in registry is updated by querying the registry state:

<!-- markdownlint-disable line-length -->
```
oasis-node registry entity list -a $ADDR -v

{"v":1,"id":"JTUtHd4XYQjh//e6eYU7Pa/XMFG88WE+jixvceIfWrk=","nodes":["LQu4ZtFg8OJ0MC4M4QMeUR7Is6Xt4A/CW+PK/7TPiH0="]}
{"v":1,"id":"+MJpnSTzc11dNI5emMa+asCJH5cxBiBCcpbYE4XBdso=","nodes":["vWUfSmjrHSlN5tSSO3/Qynzx+R/UlwPV9u+lnodQ00c="]}
{"v":1,"id":"TqUyj5Q+9vZtqu10yw6Zw7HEX3Ywe0JQA9vHyzY47TU=","allow_entity_signed_nodes":true}
```
<!-- markdownlint-enable line-length -->

Node is now able to register and the runtime should get unsuspended, make sure
this happens by querying the registry for runtimes:

```
oasis-node registry runtime list -a $ADDR -v
```

## Testing the Runtime

Now that the runtime node is running, is registered, and runtime is unsuspended,
we can test the runtime by submitting runtime transactions.
For that we use the [simple-keyvalue-client] binary which tests the
functionality of the `simple-keyvalue` runtime.

If you followed [build instructions] the built client binary is available at
target/default/debug/simple-keyvalue-client

```
./target/default/debug/simple-keyvalue-client \
  --node-address unix:/tmp/runtime-example/runtime-node/internal.sock \
  --runtime-id $RUNTIME_ID
```

[simple-keyvalue-client]: ../../tests/clients/simple-keyvalue
