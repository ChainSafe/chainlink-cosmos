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

### ~~~ BEGIN ACCOUNT TESTS ~~~ ###

# aDd AlIcE aS cHaInLiNk oRaClE iN aCcOuNt StOrE
echo "adding alice chainlink account"
addChainlinkAccountTx=$(chainlinkd tx chainlink add-chainlink-account "aliceChainlinkPubKey" "aliceChainlinkSigningKey" --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
sleep 1
addChainlinkAccountTxResp=$(echo ${addChainlinkAccountTx#*\]} | jq '.height')
if [ "$addChainlinkAccountTxResp" == "\"0\"" ]
then
  errorAndExit "Error in adding alice's chainlink account: $addChainlinkAccountTx"
fi
echo "added alice account successfully..."

# gEt AlIcE cHaInLiNk AcCoUnT iNfO
echo "getting alice chainlink account information"
getAliceAccountInfo=$(chainlinkd query chainlink getAccountInfo $(chainlinkd keys show alice -a) --from alice --keyring-backend test --chain-id testchain)

aliceSubmitterAddress=$(echo "$getAliceAccountInfo" | jq '.account.submitter')
if [ "$aliceSubmitterAddress" != "\"$aliceAddr\"" ]
then
  errorAndExit "Error incorrect account submitter address, expecting: $aliceAddr, got: $aliceSubmitterAddress"
fi
 
aliceChainlinkPublicKey=$(echo "$getAliceAccountInfo" | jq '.account.chainlinkPublicKey')
if [ "$aliceChainlinkPublicKey" != "\"YWxpY2VDaGFpbmxpbmtQdWJLZXk=\"" ]
then
  errorAndExit "Error incorrect account public key, expecting: YWxpY2VDaGFpbmxpbmtQdWJLZXk=, got: $aliceChainlinkPublicKey"
fi

aliceChainlinkSigningKey=$(echo "$getAliceAccountInfo" | jq '.account.chainlinkSigningKey')
if [ "$aliceChainlinkSigningKey" != "\"YWxpY2VDaGFpbmxpbmtTaWduaW5nS2V5\"" ]
then
  errorAndExit "Error incorrect account signing key, expecting: YWxpY2VDaGFpbmxpbmtTaWduaW5nS2V5, got: $aliceChainlinkSigningKey"
fi

alicePiggyAddress=$(echo "$getAliceAccountInfo" | jq '.account.piggyAddress')
if [ "$alicePiggyAddress" != "\"$aliceAddr\"" ]
then
  errorAndExit "Error incorrect account piggy address, expecting: $aliceAddr, got: $alicePiggyAddress"
fi
echo "got alice account info successfully..."

# mAkE sUrE bOb AcCoUnT dOeS nOt ExIsT
echo "getting bob chainlink account information"
getBobAccountInfo=$(chainlinkd query chainlink getAccountInfo $(chainlinkd keys show bob -a) --from bob --keyring-backend test --chain-id testchain)
getBobAccountInfoResp=$(echo ${getBobAccountInfo#*\]} | jq '.account.submitter')
if [ "$getBobAccountInfoResp" != "\"\"" ]
then
  errorAndExit "Error getting bob's chainlink account: $getBobAccountInfo"
fi
echo "empty account found succesfully..."

# dIsAlLoW rEpEaT aCcOuNt CrEaTiOn
echo "attempting to add alice chainlink account again"
addChainlinkAccountTx=$(chainlinkd tx chainlink add-chainlink-account "aliceChainlinkPubKey" "aliceChainlinkSigningKey" --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
sleep 1
addChainlinkAccountTxResp=$(echo ${addChainlinkAccountTx#*\]} | jq '.height')
if [ "$addChainlinkAccountTxResp" != "\"0\"" ]
then
  errorAndExit "Error in rejecting alice's chainlink account: $addChainlinkAccountTx"
fi
echo "blocked repeat account successfully..."

# eDiT aLiCe PiGgY aDdReSs tO bOb'S aDdReSs
echo "edit alice piggy address"
editAlicePiggyAddressTx=$(chainlinkd tx chainlink edit-piggy-address $(chainlinkd keys show bob -a) --from alice --keyring-backend test --chain-id testchain <<< 'y\n')

# gEt AlIcE cHaInLiNk AcCoUnT iNfO
echo "getting alice chainlink account information"
getAliceAccountInfo=$(chainlinkd query chainlink getAccountInfo $(chainlinkd keys show alice -a) --from alice --keyring-backend test --chain-id testchain)

aliceSubmitterAddress=$(echo "$getAliceAccountInfo" | jq '.account.submitter')
if [ "$aliceSubmitterAddress" != "\"$aliceAddr\"" ]
then
  errorAndExit "Error alice submitter should not have changed, expecting: $aliceAddr, got: $aliceSubmitterAddress"
fi
 
aliceChainlinkPublicKey=$(echo "$getAliceAccountInfo" | jq '.account.chainlinkPublicKey')
if [ "$aliceChainlinkPublicKey" != "\"YWxpY2VDaGFpbmxpbmtQdWJLZXk=\"" ]
then
  errorAndExit "Error alice public key should not have changed, expecting: YWxpY2VDaGFpbmxpbmtQdWJLZXk=, got: $aliceChainlinkPublicKey"
fi

aliceChainlinkSigningKey=$(echo "$getAliceAccountInfo" | jq '.account.chainlinkSigningKey')
if [ "$aliceChainlinkSigningKey" != "\"YWxpY2VDaGFpbmxpbmtTaWduaW5nS2V5\"" ]
then
  errorAndExit "Error alice signing key should not have changed, expecting: YWxpY2VDaGFpbmxpbmtTaWduaW5nS2V5, got: $aliceChainlinkSigningKey"
fi

alicePiggyAddress=$(echo "$getAliceAccountInfo" | jq '.account.piggyAddress')
if [ "$alicePiggyAddress" != "\"$bobAddr\"" ]
then
  errorAndExit "Error incorrect account piggy address, expecting: $bobAddr, got: $alicePiggyAddress"
fi
echo "edited alice piggy address passed successfully..."

pkill chainlinkd
echo "Chainlink module FEED test has exited successfully."
exit 0
