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

# aDd NeW fEeD bY aLiCe
# wIlL aDd AlIcE aDdReSs AnD pUbLiC kEy As fEeD oWnEr
echo "adding new feed by alice"
addFeedTx=$($chainlinkCMD addFeed feedid1 $aliceAddr 1 2 3 4 $aliceAddr,$alicePK --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
addFeedTxResp=$(echo "$addFeedTx" | jq '.logs')
if [ ${#addFeedTxResp} == 2 ] # log: [] if tx failed
then
  errorAndExit "Error in goodTx1: $addFeedTx"
fi

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
echo "Chainlink module EDITFEED test has exited successfully."
exit 0


