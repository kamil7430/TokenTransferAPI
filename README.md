# Token Transfer API

A simple GraphQL API backend for transferring BTP tokens between wallets.

Initially, there is only one wallet (address `0x0000000000000000000000000000000000000000`) holding 1,000,000 BTP tokens. The API supports transferring tokens from one wallet to another.

## Usage

### Prerequisites

To run the project, you need Docker and Go. 

### Running

1. Create a file `db/password.txt` and put there a PostgreSQL database's password of your choice. You can use the example file provided just by renaming it.

2. Use the following command: 

```bash
docker compose up --build
```

3. The service should be available at http://localhost:8080/

### Tests

The tests require Docker running. You can run tests using the following command:

```bash
go test ./...
```

If the tests are stuck on container creation, try to pull the image by yourself:

```bash
docker pull postgres:16-alpine
```

## GraphQL operations

### Queries

```graphql
wallet(address: Address!): Wallet!
```

Fetches the wallet with the specified address.

### Mutations

```graphql
transfer(from_address: Address!, to_address: Address!, amount: Int64!): Int64!
```

Concurrent-safe mutation that transfers `amount` tokens from wallet with `from_address` address to wallet with `to_address` address.

### Examples

```graphql
mutation {
    # Transfer some tokens
    transfer(from_address: "0x0000000000000000000000000000000000000000", to_address: "0x0000000000000000000000000000000000000001", amount: 200)
}
```

```graphql
query {
    # Query a wallet
    wallet(address: "0x0000000000000000000000000000000000000001") {
        address
        tokens
    }
}
```