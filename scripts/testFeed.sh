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

# aDd NeW fEeD bY aLiCe
# wIlL aDd AlIcE aDdReSs AnD pUbLiC kEy
echo "adding new feed by alice"
addFeedTx=$($chainlinkCMD addFeed feedid1 "this is the test feed 1" $aliceAddr 1 2 3 4 $aliceAddr,$alicePK --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
addFeedTxResp=$(echo "$addFeedTx" | jq '.logs')
if [ ${#addFeedTxResp} == 2 ] # log: [] if tx failed
then
  errorAndExit "Error in goodTx1: $addFeedTx"
fi

# iNiTiAl BaLaNcE oF aLiCe b4 rEwArD
aliceCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
if [ "$aliceCurrBal" != "\"1000000\"" ]
then
  errorAndExit "Error in initial distribution; expected 1000000, got $aliceCurrBal"
fi

# iNiTiAl BaLaNcE oF bOb B4 rEwArD
bobCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show bob -a) --denom link --output json | jq '.amount')
if [ "$bobCurrBal" != "\"1000000\"" ]
then
  errorAndExit "Error in initial distribution; expected \"1000000\", got $bobCurrBal"
fi

# sUbMiT fEeD dAtA bY aLiCe
echo "submitting feed data by alice"
submitFeedTx1=$($chainlinkCMD submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
submitFeedTx1Resp=$(echo "$submitFeedTx1" | jq '.height')
if [ "$submitFeedTx1Resp" == "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #1: $submitFeedTx1"
fi

# cHeCk If AlIcE gOt ThE rEwArD
echo "checking alice's reward distribution #1"
aliceCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
if [ "$aliceCurrBal" != "\"1000004\"" ]
then
  errorAndExit "Error in reward distribution for alice; expected \"1000004\", got $aliceCurrBal"
fi

# bOb ShOuLd NoT aNy rEwArD
echo "checking bob's reward distribution #1"
bobCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show bob -a) --denom link --output json | jq '.amount')
if [ "$bobCurrBal" != "\"1000000\"" ]
then
  errorAndExit "Error in reward distribution; expected \"1000000\", got $bobCurrBal"
fi

# sUbMiT fEeD dAtA bY cErLo (nOn-AuThOrIzEd DaTa PrOvIdEr)...
echo "submitting feed data by unauthorized data provider"
badSubmitFeedTx=$($chainlinkCMD submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from bob --keyring-backend test --chain-id testchain <<< 'y\n')
badSubmitFeedTxResp=$(echo "$badSubmitFeedTx" | jq '.raw_log')
if [ "$badSubmitFeedTxResp" != "\"invalid data provider: unauthorized\"" ]
then
  errorAndExit "Error in sending bad feed data: $badSubmitFeedTx"
fi

##############

# aDd BoB aS dAtA pRoViDeR
echo "adding bob as a data provider"
addBobTx=$($chainlinkCMD addDataProvider feedid1 $bobAddr $bobPK --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
addBobTxResp=$(echo $addBobTx | jq '.height')
if [ "$addBobTxResp" == "\"0\"" ]
then
  errorAndExit "Error in adding bob as a data provider: $addBobTx"
fi

# uPdAtE fEeD rEwArD
echo "updating feed reward to $newFeedReward"
newFeedReward=100
updateFeedReward=$($chainlinkCMD setFeedReward feedid1 $newFeedReward --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
updateFeedRewardResp=$(echo "$updateFeedReward" | jq '.height')
if [ "$updateFeedRewardResp" == "\"0\"" ]
then
  errorAndExit "Error in updating feed reward: $updateFeedReward"
fi

# sUbMiT fEeD dAtA bY bOb
echo "submitting feed data by bob"
submitFeedTx2=$($chainlinkCMD submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from bob --keyring-backend test --chain-id testchain <<< 'y\n')
submitFeedTx2Resp=$(echo "$submitFeedTx2" | jq '.height')
if [ "$submitFeedTx2Resp" == "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #2: $submitFeedTx2"
fi

# cHeCk If AlIcE gOt ThE uPdAtEd ReWaRd
echo "checking alice's reward distribution #2"
aliceCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
if [ "$aliceCurrBal" != "\"1000104\"" ]
then
  errorAndExit "Error in reward distribution; expected \"1000104\", got $aliceCurrBal"
fi

# bOb ShOuLd NoW gEt rEwArD
echo "checking bob's reward distribution #2"
bobCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show bob -a) --denom link --output json | jq '.amount')
if [ "$bobCurrBal" != "\"1000100\"" ]
then
  errorAndExit "Error in reward distribution; expected \"1000100\", got $bobCurrBal"
fi

### ~~~ BEGIN FEED EDIT TESTS ~~~ ###

# fEeD oWnEr WiLl eDiT
# eDiT sUbMiSsIoN cOuNt 
echo "edit feed's submission count by alice"
setSubmissionCountTx=$($chainlinkCMD setSubmissionCount feedid1 10 --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
setSubmissionCountTxResp=$(echo "$setSubmissionCountTx" | jq '.height')
if [ "$setSubmissionCountTxResp" == "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #1: $setSubmissionCountTx"
fi

# eDiT hEaRtBeAt TrIgGeR
echo "edit feed's heartbeat trigger by alice"
setHeartbeatTriggerTx=$($chainlinkCMD setHeartbeatTrigger feedid1 20 --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
setHeartbeatTriggerTxResp=$(echo "$setHeartbeatTriggerTx" | jq '.height')
if [ "$setHeartbeatTriggerTxResp" == "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #1: $setHeartbeatTriggerTx"
fi

# eDiT dEvIaTiOn ThReShOlD tRiGgEr
echo "edit feed's deviation threshold trigger by alice"
setDeviationThresholdTriggerTx=$($chainlinkCMD setDeviationThresholdTrigger feedid1 30 --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
setDeviationThresholdTriggerTxResp=$(echo "$setDeviationThresholdTriggerTx" | jq '.height')
if [ "$setHeartbeatTriggerTxResp" == "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #1: $setDeviationThresholdTriggerTx"
fi

# eDiT fEeD rEwArD
echo "edit feed's reward by alice"
setFeedRewardTx=$($chainlinkCMD setFeedReward feedid1 40 --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
setFeedRewardTxResp=$(echo "$setFeedRewardTx" | jq '.height')
if [ "$setFeedRewardTxResp" == "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #1: $setFeedRewardTx"
fi

# nOn-FeEd OwMeR cAnNoT eDiT
# eDiT sUbMiSsIoN cOuNt 
echo "edit feed's submission count by cerlo"
setSubmissionCountTx=$($chainlinkCMD setSubmissionCount feedid1 10 --from cerlo --keyring-backend test --chain-id testchain <<< 'y\n')
setSubmissionCountTxResp=$(echo "$setSubmissionCountTx" | jq '.height')
if [ "$setSubmissionCountTxResp" != "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #1: $setSubmissionCountTx"
fi

# eDiT hEaRtBeAt TrIgGeR
echo "edit feed's heartbeat trigger by cerlo"
setHeartbeatTriggerTx=$($chainlinkCMD setHeartbeatTrigger feedid1 20 --from cerlo --keyring-backend test --chain-id testchain <<< 'y\n')
setHeartbeatTriggerTxResp=$(echo "$setHeartbeatTriggerTx" | jq '.height')
if [ "$setHeartbeatTriggerTxResp" != "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #1: $setHeartbeatTriggerTx"
fi

# eDiT dEvIaTiOn ThReShOlD tRiGgEr
echo "edit feed's deviation threshold trigger by cerlo"
setDeviationThresholdTriggerTx=$($chainlinkCMD setDeviationThresholdTrigger feedid1 30 --from cerlo --keyring-backend test --chain-id testchain <<< 'y\n')
setDeviationThresholdTriggerTxResp=$(echo "$setDeviationThresholdTriggerTx" | jq '.height')
if [ "$setDeviationThresholdTriggerTxResp" != "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #1: $setDeviationThresholdTriggerTx"
fi

# eDiT fEeD rEwArD
echo "edit feed's reward by cerlo"
setFeedRewardTx=$($chainlinkCMD setFeedReward feedid1 40 --from cerlo --keyring-backend test --chain-id testchain <<< 'y\n')
setFeedRewardTxResp=$(echo "$setFeedRewardTx" | jq '.height')
if [ "$setFeedRewardTxResp" != "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #1: $setFeedRewardTx"
fi


pkill chainlinkd
echo "Chainlink module FEED test has exited successfully."
exit 0