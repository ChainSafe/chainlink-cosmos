#!/bin/bash

#### according to `start.sh`, ALICE is the Module Owner and will add BOB as a Feed Owner. #####
./scripts/start.sh > "$(pwd)"/chainlinkd.log 2>&1 &
sleep 10

aliceAddr=$(chainlinkd keys show alice -a)
alicePK=$(chainlinkd keys show alice -p)

bobAddr=$(chainlinkd keys show bob -a)
bobPK=$(chainlinkd keys show bob -p)

cerloAddr=$(chainlinkd keys show cerlo -a)
cerloPK=$(chainlinkd keys show cerlo -p)

# aDd NeW fEeD bY aLiCe
# wIlL aDd AlIcE aDdReSs AnD pUbLiC kEy
echo "adding new feed by alice"
addFeedTx=$(chainlinkd tx chainlink addFeed feedid1 $aliceAddr 1 2 3 4 $aliceAddr,$alicePK --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
addFeedTxResp=$(echo "$addFeedTx" | jq '.raw_log')
# "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"AddFeed\"}]}]}]"
if [ "$addFeedTxResp" != "\"[{\\\"events\\\":[{\\\"type\\\":\\\"message\\\",\\\"attributes\\\":[{\\\"key\\\":\\\"action\\\",\\\"value\\\":\\\"AddFeed\\\"}]}]}]\"" ]
then
  echo "Error in goodTx1: $addFeedTx"
  pkill chainlinkd
  exit 1
fi

# iNiTiAl BaLaNcE oF aLiCe b4 rEwArD
aliceCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
if [ "$aliceCurrBal" != "\"1000000\"" ]
then
  echo "Error in initial distribution; expected 1000000, got $aliceCurrBal"
  pkill chainlinkd
  exit 1
fi

# iNiTiAl BaLaNcE oF bOb B4 rEwArD
bobCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show bob -a) --denom link --output json | jq '.amount')
if [ "$bobCurrBal" != "\"1000000\"" ]
then
  echo "Error in initial distribution; expected \"1000000\", got $bobCurrBal"
  pkill chainlinkd
  exit 1
fi

# sUbMiT fEeD dAtA bY aLiCe
echo "submitting feed data by alice"
submitFeedTx1=$(chainlinkd tx chainlink submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
submitFeedTx1Resp=$(echo "$submitFeedTx1" | jq '.height')
if [ "$submitFeedTx1Resp" == "\"0\"" ]
then
  echo "Error in submitting feed data #1: $submitFeedTx1"
  pkill chainlinkd
  exit 1
fi

# cHeCk If AlIcE gOt ThE rEwArD
echo "checking alice's reward distribution #1"
aliceCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
if [ "$aliceCurrBal" != "\"1000004\"" ]
then
  echo "Error in reward distribution for alice; expected \"1000004\", got $aliceCurrBal"
  pkill chainlinkd
  exit 1
fi

# bOb ShOuLd NoT aNy rEwArD
echo "checking bob's reward distribution #1"
bobCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show bob -a) --denom link --output json | jq '.amount')
if [ "$bobCurrBal" != "\"1000000\"" ]
then
  echo "Error in reward distribution; expected \"1000000\", got $bobCurrBal"
  pkill chainlinkd
  exit 1
fi

# sUbMiT fEeD dAtA bY cErLo (nOn-AuThOrIzEd DaTa PrOvIdEr)...
echo "submitting feed data by unauthorized data provider"
badSubmitFeedTx=$(chainlinkd tx chainlink submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from bob --keyring-backend test --chain-id testchain <<< 'y\n')
badSubmitFeedTxResp=$(echo "$badSubmitFeedTx" | jq '.raw_log')
if [ "$badSubmitFeedTxResp" != "\"invalid data provider: unauthorized\"" ]
then
  echo "Error in sending bad feed data: $badSubmitFeedTx"
  pkill chainlinkd
  exit 1
fi

##############

# aDd BoB aS dAtA pRoViDeR
echo "adding bob as a data provider"
addBobTx=$(chainlinkd tx chainlink addDataProvider feedid1 $bobAddr $bobPK --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
addBobTxResp=$(echo $addBobTx | jq '.height')
if [ "$addBobTxResp" == "\"0\"" ]
then
  echo "Error in adding bob as a data provider: $addBobTx"
  pkill chainlinkd
  exit 1
fi

# uPdAtE fEeD rEwArD
echo "updating feed reward to $newFeedReward"
newFeedReward=100
updateFeedReward=$(chainlinkd tx chainlink setFeedReward feedid1 $newFeedReward --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
updateFeedRewardResp=$(echo "$updateFeedReward" | jq '.height')
if [ "$updateFeedRewardResp" == "\"0\"" ]
then
  echo "Error in updating feed reward: $updateFeedReward"
  pkill chainlinkd
  exit 1
fi

# sUbMiT fEeD dAtA bY bOb
echo "submitting feed data by bob"
submitFeedTx2=$(chainlinkd tx chainlink submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from bob --keyring-backend test --chain-id testchain <<< 'y\n')
submitFeedTx2Resp=$(echo "$submitFeedTx2" | jq '.height')
if [ "$submitFeedTx2Resp" == "\"0\"" ]
then
  echo "Error in submitting feed data #2: $submitFeedTx2"
  pkill chainlinkd
  exit 1
fi

# cHeCk If AlIcE gOt ThE uPdAtEd ReWaRd
echo "checking alice's reward distribution #2"
aliceCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
if [ "$aliceCurrBal" != "\"1000104\"" ]
then
  echo "Error in reward distribution; expected \"1000104\", got $aliceCurrBal"
  pkill chainlinkd
  exit 1
fi

# bOb ShOuLd NoW gEt rEwArD
echo "checking bob's reward distribution #2"
bobCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show bob -a) --denom link --output json | jq '.amount')
if [ "$bobCurrBal" != "\"1000100\"" ]
then
  echo "Error in reward distribution; expected \"1000100\", got $bobCurrBal"
  pkill chainlinkd
  exit 1
fi

pkill chainlinkd
echo "Chainlink module ADDFEED test has exited successfully."
exit 0
