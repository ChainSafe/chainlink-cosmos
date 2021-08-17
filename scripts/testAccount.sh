#!/bin/bash

function errorAndExit() {
  echo $1
  pkill chainlinkd
  exit 1
}

chainlinkCMD="chainlinkd tx chainlink"

#### according to `start.sh`, ALICE is the Module Owner. #####
./scripts/start.sh > "$(pwd)"/chainlinkd.log 2>&1 &
sleep 10

aliceAddr=$(chainlinkd keys show alice -a)
alicePK=$(chainlinkd keys show alice -p)

bobAddr=$(chainlinkd keys show bob -a)
bobPK=$(chainlinkd keys show bob -p)

cerloAddr=$(chainlinkd keys show cerlo -a)
cerloPK=$(chainlinkd keys show cerlo -p)

### ~~~ BEGIN FEED ADD TESTS ~~~ ###

# aDd AlIcE aS cHaInLiNk oRaClE iN aCcOuNt StOrE
echo "adding alice chainlink account"
addChainlinkAccountTx=$(chainlinkd tx chainlink add-chainlink-account "aliceChainlinkPubKey" "aliceChainlinkSigningKey" --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
sleep 1
addChainlinkAccountTxResp=$(echo ${addChainlinkAccountTx#*\]} | jq '.height')
if [ "$addChainlinkAccountTxResp" == "\"0\"" ]
then
  errorAndExit "Error in adding alice's chainlink account: $addChainlinkAccountTx"
fi

# gEt AlIcE cHaInLiNk AcCoUnT iNfO
echo "getting alice chainlink account information"
getAliceAccountInfo=$(chainlinkd query chainlink getAccountInfo $(chainlinkd keys show alice -a) --from alice --keyring-backend test --chain-id testchain)

aliceSubmitterAddress=$(echo "$getAliceAccountInfo" | jq '.account.submitter')
echo "$aliceSubmitterAddress"
echo "$aliceAddr"
if [ "$aliceSubmitterAddress" != "\"$aliceAddr\"" ]
then
  errorAndExit "Error incorrect account submitter address, expecting: $aliceAddr, got: $aliceSubmitterAddress"
fi
 
aliceChainlinkPublicKey=$(echo "$getAliceAccountInfo" | jq '.account.chainlinkPublicKey')
echo $aliceChainlinkPublicKey
if [ "$aliceChainlinkPublicKey" != "\"YWxpY2VDaGFpbmxpbmtQdWJLZXk=\"" ]
then
  errorAndExit "Error incorrect account public key, expecting: YWxpY2VDaGFpbmxpbmtQdWJLZXk=, got: $aliceChainlinkPublicKey"
fi

aliceChainlinkSigningKey=$(echo "$getAliceAccountInfo" | jq '.account.chainlinkSigningKey')
echo $aliceChainlinkSigningKey
if [ "$aliceChainlinkSigningKey" != "\"YWxpY2VDaGFpbmxpbmtTaWduaW5nS2V5\"" ]
then
  errorAndExit "Error incorrect account signing key, expecting: YWxpY2VDaGFpbmxpbmtTaWduaW5nS2V5, got: $aliceChainlinkSigningKey"
fi

alicePiggyAddress=$(echo "$getAliceAccountInfo" | jq '.account.piggyAddress')
echo $alicePiggyAddress
if [ "$alicePiggyAddress" != "\"$aliceAddr\"" ]
then
  errorAndExit "Error incorrect account piggy address, expecting: $aliceAddr, got: $alicePiggyAddress"
fi

pkill chainlinkd
echo "Chainlink module FEED test has exited successfully."
exit 0
