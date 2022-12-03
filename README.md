# Livepool Service

## Authentication

### HMAC Authentication

Used by livepool servers

### ECDSA Authentication

Used by transcoders to execute payouts or retrieve information on their transcoder, jobs and nodes.

Authentication is done by having the user generate an ECDSA signature over the hash of the request body using the private key corresponding to the transcoder's registered Ethereum address.

Might switch to Sign-In-With-Ethereum: <https://docs.login.xyz/libraries/go>

## Routes

### transcoders

#### GET

Get all transcoders, a single transcoder or transcoders per region.
Displays only public data.

#### POST

Requires HMAC authentication, used by livepool servers to register a new transcoder

#### PUT

requires ECDSA authentication, meant to change a transcoder's ETH address by the transcoder itself

### nodes

#### GET

requires ECDSA authentication, gets the nodes for a transcoder.

#### POST

requires HMAC authentication, used by livepool servers to register a new node to a transcoder

### jobs

#### GET

requires ECDSA authentication, retrieves jobs for a transcoder.

#### POST

Requires HMAC authentication, can only be used by livepool servers.
Inserts a new job into the database and updates the transcoder's pending balance.
