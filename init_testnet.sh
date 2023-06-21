#{
#   "address": "9472F66238D1C730BA4081DE1B10362EF4634AE4",
#   "pub_key": {
#     "type": "tendermint/PubKeyEd25519",
#     "value": "Q3tNXrIRokD05nNdNFZSROSU81kI47BbnL7zkcxB/WY="
#   },
#   "priv_key": {
#     "type": "tendermint/PrivKeyEd25519",
#     "value": "jGIbnQEh+vfC6kLKAO9fceulfV3g6cSLlt0SUZEKMKlDe01eshGiQPTmc100VlJE5JTzWQjjsFucvvORzEH9Zg=="
#   }
#}
WOOD="wood"
WOOD_PUBKEY='{"@type":"/cosmos.crypto.ed25519.PubKey","key":"Q3tNXrIRokD05nNdNFZSROSU81kI47BbnL7zkcxB/WY="}'
ARHEN="arhen"
#ARHEN_PUBKEY='{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A+RwzqKr4NyhIsAK7NeXeFwbZm5vA6Qt0A1rh7YqX3lo"}'
ARHEN_PUBKEY='{"@type":"/cosmos.crypto.ed25519.PubKey","key":"3N7y8v7SMkwgRoOhtyZ6/WHR2qcFjyFV2hhHCrZhShE="}'
GAVIN="gavin"
#GAVIN_PUBKEY='{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A/JbNhr249zvdIL1UaPfG2dESBug2fcq4rRK179aTszy"}'
GAVIN_PUBKEY='{"@type":"/cosmos.crypto.ed25519.PubKey","key":"u8tfMH1+46ZZOV7ffkFezYUdGWN3+YxEuq96e6PQ8UI="}'
CHAINID="canto_7701-1"
MONIKER="plex-validator"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# to trace evm
#TRACE="--trace"
TRACE=""

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# Reinstall daemon
rm -rf ~/.cantod*
make install

# Set client config
cantod config keyring-backend $KEYRING
cantod config chain-id $CHAINID

# if $KEY exists it should be deleted
cantod keys add $WOOD --keyring-backend $KEYRING --algo $KEYALGO
cantod keys add $ARHEN --keyring-backend $KEYRING --algo $KEYALGO
cantod keys add $GAVIN --keyring-backend $KEYRING --algo $KEYALGO

# Set moniker and chain-id for Canto (Moniker can be anything, chain-id must be an integer)
cantod init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to acanto
cat $HOME/.cantod/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="acanto"' > $HOME/.cantod/config/tmp_genesis.json && mv $HOME/.cantod/config/tmp_genesis.json $HOME/.cantod/config/genesis.json
cat $HOME/.cantod/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="acanto"' > $HOME/.cantod/config/tmp_genesis.json && mv $HOME/.cantod/config/tmp_genesis.json $HOME/.cantod/config/genesis.json
cat $HOME/.cantod/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="acanto"' > $HOME/.cantod/config/tmp_genesis.json && mv $HOME/.cantod/config/tmp_genesis.json $HOME/.cantod/config/genesis.json
cat $HOME/.cantod/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="acanto"' > $HOME/.cantod/config/tmp_genesis.json && mv $HOME/.cantod/config/tmp_genesis.json $HOME/.cantod/config/genesis.json
cat $HOME/.cantod/config/genesis.json | jq '.app_state["inflation"]["params"]["mint_denom"]="acanto"' > $HOME/.cantod/config/tmp_genesis.json && mv $HOME/.cantod/config/tmp_genesis.json $HOME/.cantod/config/genesis.json

# Change voting params so that submitted proposals pass immediately for testing
cat $HOME/.cantod/config/genesis.json| jq '.app_state.gov.voting_params.voting_period="30s"' > $HOME/.cantod/config/tmp_genesis.json && mv $HOME/.cantod/config/tmp_genesis.json $HOME/.cantod/config/genesis.json


# disable produce empty block
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.cantod/config/config.toml
  else
    sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.cantod/config/config.toml
fi

if [[ $1 == "pending" ]]; then
  if [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.cantod/config/config.toml
      sed -i '' 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.cantod/config/config.toml
      sed -i '' 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.cantod/config/config.toml
      sed -i '' 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.cantod/config/config.toml
      sed -i '' 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.cantod/config/config.toml
      sed -i '' 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.cantod/config/config.toml
      sed -i '' 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.cantod/config/config.toml
      sed -i '' 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.cantod/config/config.toml
      sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.cantod/config/config.toml
  else
      sed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.cantod/config/config.toml
      sed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.cantod/config/config.toml
      sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.cantod/config/config.toml
      sed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.cantod/config/config.toml
      sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.cantod/config/config.toml
      sed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.cantod/config/config.toml
      sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.cantod/config/config.toml
      sed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.cantod/config/config.toml
      sed -i 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.cantod/config/config.toml
  fi
fi

# Allocate genesis accounts (cosmos formatted addresses)
cantod add-genesis-account $WOOD 964723926400000000000000000acanto --keyring-backend $KEYRING
cantod add-genesis-account $ARHEN 35276073600000000000000000acanto --keyring-backend $KEYRING
cantod add-genesis-account $GAVIN 50000000000000000000000000acanto --keyring-backend $KEYRING

# Update total supply with claim values
#validators_supply=$(cat $HOME/.cantod/config/genesis.json | jq -r '.app_state["bank"]["supply"][0]["amount"]')
# Bc is required to add this big numbers
# total_supply=$(bc <<< "$amount_to_claim+$validators_supply")
total_supply=1050000000000000000000000000
cat $HOME/.cantod/config/genesis.json | jq -r --arg total_supply "$total_supply" '.app_state["bank"]["supply"][0]["amount"]=$total_supply' > $HOME/.cantod/config/tmp_genesis.json && mv $HOME/.cantod/config/tmp_genesis.json $HOME/.cantod/config/genesis.json

echo $KEYRING
echo $KEY
# Sign genesis transaction
mkdir $HOME/.cantod/config/gentx
cantod gentx $WOOD 900000000000000000000000acanto --keyring-backend $KEYRING --chain-id $CHAINID --output-document $HOME/.cantod/config/gentx/gentx-1.json --pubkey $WOOD_PUBKEY
#cantod gentx $ARHEN 100000000000000000000000acanto --keyring-backend $KEYRING --chain-id $CHAINID --output-document $HOME/.cantod/config/gentx/gentx-2.json --pubkey $ARHEN_PUBKEY
#cantod gentx $GAVIN 100000000000000000000000acanto --keyring-backend $KEYRING --chain-id $CHAINID --output-document $HOME/.cantod/config/gentx/gentx-3.json --pubkey $GAVIN_PUBKEY
#cantod gentx $KEY2 1000000000000000000000acanto --keyring-backend $KEYRING --chain-id $CHAINID

# Collect genesis tx
cantod collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
cantod validate-genesis

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
#cantod start --pruning=nothing --trace --log_level trace --minimum-gas-prices=1.000acanto --json-rpc.api eth,txpool,personal,net,debug,web3 --rpc.laddr "tcp://0.0.0.0:26657" --api.enable true --api.enabled-unsafe-cors true

