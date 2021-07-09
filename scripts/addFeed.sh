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
goodTx1=$(chainlinkd tx chainlink addFeed feedid1 $aliceAddr 1 2 3 $aliceAddr,$alicePK --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
goodTx1Resp=$(echo "$goodTx1" | jq '.raw_log')
# "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"AddFeed\"}]}]}]"
echo "sending goodTx1"
if [ "$goodTx1Resp" != "\"[{\\\"events\\\":[{\\\"type\\\":\\\"message\\\",\\\"attributes\\\":[{\\\"key\\\":\\\"action\\\",\\\"value\\\":\\\"AddFeed\\\"}]}]}]\"" ]

then
  echo "Error in goodTx1: $goodTx1Resp"
  pkill chainlinkd
  exit 1
fi

# sUbMiT fEeD dAtA bY aLiCe
goodTx2=$(chainlinkd tx chainlink submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
goodTx2Resp=$(echo "$goodTx2" | jq '.raw_log')
echo "sending goodTx2"
# "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"SubmitFeedData\"}]}]}]"
if [ "$goodTx2Resp" != "\"[{\\\"events\\\":[{\\\"type\\\":\\\"message\\\",\\\"attributes\\\":[{\\\"key\\\":\\\"action\\\",\\\"value\\\":\\\"SubmitFeedData\\\"}]}]}]\"" ]
then
  echo "Error in goodTx2: $goodTx2Resp"
  pkill chainlinkd
  exit 1
fi

# sUbMiT fEeD dAtA bY bOb (nOn-AuThOrIzEd DaTa PrOvIdEr)...
badTx1=$(chainlinkd tx chainlink submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from bob --keyring-backend test --chain-id testchain <<< 'y\n')
badTx1Resp=$(echo "$badTx1" | jq '.raw_log')
# "raw_log":"invalid data provider: unauthorized"
echo "sending badTx1"
if [ "$badTx1Resp" != "\"invalid data provider: unauthorized\"" ]
then
  echo "Error in badTx1: $badTx1Resp"
  pkill chainlinkd
  exit 1
fi
echo "badTx1 rejected successfully"

# aDd NeW dAtA pRoViDeR


pkill chainlinkd
echo "Chainlink module ADDFEED test has exited successfully."
exit 0