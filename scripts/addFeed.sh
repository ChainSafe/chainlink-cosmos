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
# wIlL uSe AlIcE aDdReSs AnD pUbLiC kEy
addFeedTx=$(chainlinkd tx chainlink addFeed feedid1 $aliceAddr 1 2 3 4 $aliceAddr,$alicePK --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
addFeedTxResp=$(echo "$addFeedTx" | jq '.raw_log')
# "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"AddFeed\"}]}]}]"
echo "sending goodTx1"
if [ "$addFeedTxResp" != "\"[{\\\"events\\\":[{\\\"type\\\":\\\"message\\\",\\\"attributes\\\":[{\\\"key\\\":\\\"action\\\",\\\"value\\\":\\\"AddFeed\\\"}]}]}]\"" ]
then
  echo "Error in goodTx1: $addFeedTxResp"
  pkill chainlinkd
  exit 1
fi

# iNiTiAl BaLaNcE oF aLiCe b4 rEwArD
aliceInitBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
if [ "$aliceInitBal" != "\"1000000\"" ]
then
  echo "Error in initial distribution; expected 1000000, got $aliceCurrBal"
  pkill chainlinkd
  exit 1
fi

# sUbMiT fEeD dAtA bY aLiCe
submitFeedTx=$(chainlinkd tx chainlink submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
submitFeedTxResp=$(echo "$submitFeedTx" | jq '.height')
echo "sending goodTx2"
if [ "$submitFeedTxResp" == "\"0\"" ]
then
  echo "Error in goodTx2: $submitFeedTxResp"
  pkill chainlinkd
  exit 1
fi

# cHeCk If AlIcE gOt ThE rEwArD
aliceCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
echo "checking reward distribution"
if [ "$aliceCurrBal" != "\"1000004\"" ]
then
  echo "Error in reward distribution; expected 1000004, got $aliceCurrBal"
  pkill chainlinkd
  exit 1
fi

# sUbMiT fEeD dAtA bY bOb (nOn-AuThOrIzEd DaTa PrOvIdEr)...
badSubmitFeedTx=$(chainlinkd tx chainlink submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from bob --keyring-backend test --chain-id testchain <<< 'y\n')
badSubmitFeedTxResp=$(echo "$badSubmitFeedTx" | jq '.raw_log')
# "raw_log":"invalid data provider: unauthorized"
echo "sending badTx1"
if [ "$badSubmitFeedTxResp" != "\"invalid data provider: unauthorized\"" ]
then
  echo "Error in badTx1: $badSubmitFeedTxResp"
  pkill chainlinkd
  exit 1
fi
echo "badTx1 rejected successfully"

##############

# uPdAtE fEeD rEwArD
newFeedReward=100
updateFeedReward=$(chainlinkd tx chainlink setFeedReward feedid1 $newFeedReward --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
updateFeedRewardResp=$(echo "$updateFeedReward" | jq '.height')
echo "updaing feed reward to $newFeedReward"
if [ "$updateFeedRewardResp" == "\"0\"" ]
then
  echo "Error in goodTx2: $updateFeedRewardResp"
  pkill chainlinkd
  exit 1
fi

# sUbMiT fEeD dAtA bY aLiCe
submitFeedTx2=$(chainlinkd tx chainlink submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
submitFeedTx2Resp=$(echo "$submitFeedTx2" | jq '.height')
echo "sending goodTx2"
if [ "$submitFeedTx2Resp" == "\"0\"" ]
then
  echo "Error in goodTx2: $submitFeedTx2Resp"
  pkill chainlinkd
  exit 1
fi

# cHeCk If AlIcE gOt ThE uPdAtEd ReWaRd
aliceCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
echo "checking reward distribution"
if [ "$aliceCurrBal" != "\"1000104\"" ]
then
  echo "Error in reward distribution; expected 1000104, got $aliceCurrBal"
  pkill chainlinkd
  exit 1
fi

pkill chainlinkd
echo "Chainlink module ADDFEED test has exited successfully."
exit 0
