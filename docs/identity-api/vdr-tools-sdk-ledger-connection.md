# Ledger connections in VDR Tools SDK

## Overview

This page describes how [Evernym VDR Tools](https://gitlab.com/evernym/verity/vdr-tools) connects to the cheqd network ledger "pool".

It is worth noting here that the terminology of "pool" connection is specifically a legacy term originally used in [Hyperledger Indy](https://github.com/hyperledger/indy-node), which as a permissioned blockchain assumes there is a finite pools of servers. While this paradigm is no longer true in the public, permissionless world of the cheqd network, the identity APIs in VDR Tools SDK and similar Hyperledger Aries-based frameworks is retained for explanations.

## Ledger pool connection methods

Establishing a ledger "pool" connection in VDR Tools SDK broadly has the following steps:

1. Generate keys or restore them from mnemonic as described in [key management using VDR Tools SDK](vdr-tools-sdk-accounts-keys.md).
2. Add ledger "pool" configuration as described using the methods below, i.e., `indy_cheqd_pool_add`. This only adds the configuration, without actually establishing the connection. The connection is established once the first transaction is sent.

For compatibility purposes, VDR Tools SDK method names use the `indy_` prefix. This may be updated in the future as work is done on the upstream project to refactor method names to be ledger-agnostic.

### indy_cheqd_pool_add

Add a new cheqd network ledger `PoolConfig` configuration.

#### Input parameters

* `alias` (string): Friendly-name for pool connection
* `rpc_address` (string): Tendermint RPC endpoint (e.g., `http://localhost:26657`) for a cheqd network node(s) to send/receive transactions to.
* `chain_id` (string): cheqd network identifier, e.g., `cheqd-mainnet-1`

#### Example output

```jsonc
{
    "alias": "cheqd_pool",
    "rpc_address": "https://rpc.testnet.cheqd.network:443",
    "chain_id": "cheqd-testnet-4"
}
```

### indy_cheqd_pool_get_config

Fetch pool configuration for a specific connection `alias`.

#### Input parameters

* `alias` (string): Friendly-name for pool connection

#### Example output

```jsonc
{
    "alias": "cheqd_pool",
    "rpc_address": "https://rpc.testnet.cheqd.network:443",
    "chain_id": "cheqd-testnet-4"
}
```

### indy_cheqd_pool_get_all_config

Display pool configuration for all pools.

#### Example output

```js
[Object({
	"alias": String("cheqd_pool_1"),
	"chain_id": String("cheqd-testnet-4"),
	"rpc_address": String("https://rpc.testnet.cheqd.network:443")
}), 
Object({
	"alias": String("cheqd_pool_2"),
	"chain_id": String("cheqd-mainnet-1"),
	"rpc_address": String("https://rpc.cheqd.net:443")
})]
```

### indy_cheqd_pool_broadcast_tx_commit

Broadcast a signed cheqd/Cosmos transaction to node(s) in a defined pool. This wraps up any identity-related payloads generated by VDR Tools SDK in a cheqd ledger transaction wrapper [using standard Cosmos broadcast methods](https://docs.cosmos.network/master/run-node/txs.html#broadcasting-a-transaction).

#### Input parameters

* `alias` (string): Friendly-name for pool connection
* `signed_tx_raw` (string): String of bytes containing a correctly formattted cheqd/Cosmos transaction.
* `signed_tx_len` (integer): Length of signed transaction string in bytes.

#### Example output

```jsonc
{
	"check_tx": {
		"code": 0,
		"data": "",
		"log": "[]",
		"info": "",
		"gas_wanted": "300000",
		"gas_used": "38591",
		"events": [],
		"codespace": ""
	},
	"deliver_tx": {
		"code": 0,
		"data": "Cg8KCUNyZWF0ZU55bRICCAU=",
		"log": [{
			"events ": [{
				"type ": "message",
				"attributes": [{
					"key": "action",
					"value": "CreateDid"
				}]
			}]
		}],
		"info": "",
		"gas_wanted": "300000",
		"gas_used": "46474",
		"events": [{
			"type": "message",
			"attributes": [{
				"key": "YWN0aW9u",
				"value": "Q3JlYXRlTnlt"
			}]
		}],
		"codespace": ""
	},
	"hash": "364441EDC5266A0B6AF5A67D4F05AC5D1FE95BFEDFBEBBE195723BEDBA877CAE",
	"height": "121"
}
```

### indy_cheqd_pool_abci_query

Send a [Tendermint ABCI query](https://docs.cosmos.network/v0.44/intro/sdk-app-architecture.html#abci) to specified pool `alias`. ABCI queries allow custom queries to be constructed and resultant answers fetched from the cheqd network ledger, for any data that may not be covered under the usual RPC/REST API endpoints.

#### Input parameters

* `alias` (string): Friendly-name for pool connection
* `req_json` (string): ABCI query in JSON format

#### Example output

```jsonc
{
   "did":
   {
      "creator": "cheqd1x33xkjd3gqlfhz5l9h60m53pr2mdd4y3nc86h0",
      "id": 4,
      "alias": "test-alias",
      "verkey": "did:cheqd:<namespace>:<unique-id>#key1",
      "did": "did:cheqd:<namespace>:<unique-id>"
   }
}
```

### indy_cheqd_pool_abci_info

Display pool information for a specified pool `alias`. Similar to the response that can be fetched directly from Tendermint RPC endpoint at the `/abci_info` path, e.g., `http://localhost:26657/abci_info`

#### Example output

```jsonc
{
	"response": {
    	"data": "cheqd-node",
    	"version": "0.4.0",
    	"app_version": "1",
    	"last_block_height": "541557",
    	"last_block_app_hash": "fMbrqSo1KFPeKBASylG4lEVy7iItUwGVqSUw1CE9Ydw="
    }
}
```