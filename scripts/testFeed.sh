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
addFeedTx=$($chainlinkCMD add-feed feedid1 "this is the test feed 1" $aliceAddr 1 2 3 100 "" $aliceAddr,$alicePK --from alice --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
addFeedTxResp=$(echo "$addFeedTx" | jq '.logs')
if [ ${#addFeedTxResp} == 2 ] # log: [] if tx failed
then
  errorAndExit "Error in goodTx1: $addFeedTx"
fi

# iNiTiAl BaLaNcE oF aLiCe b4 rEwArD
aliceCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
if [ "$aliceCurrBal" != "\"999997\"" ]
then
  errorAndExit "Error in initial distribution; expected 999997, got $aliceCurrBal"
fi

# iNiTiAl BaLaNcE oF bOb B4 rEwArD
bobCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show bob -a) --denom link --output json | jq '.amount')
if [ "$bobCurrBal" != "\"1000000\"" ]
then
  errorAndExit "Error in initial distribution; expected \"1000000\", got $bobCurrBal"
fi

# aDd AlIcE aS cHaInLiNk oRaClE iN aCcOuNt StOrE
echo "adding alice chainlink account"
addChainlinkAccountTx=$(chainlinkd tx chainlink add-chainlink-account "aliceChainlinkPubKey" "aliceChainlinkSigningKey" --from alice --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
sleep 1
addChainlinkAccountTxResp=$(echo ${addChainlinkAccountTx#*\]} | jq '.height')
if [ "$addChainlinkAccountTxResp" == "\"0\"" ]
then
  errorAndExit "Error in adding alice's chainlink account: $addChainlinkAccountTx"
fi
echo "added alice account successfully..."

# sUbMiT fEeD dAtA bY aLiCe
echo "submitting feed data by alice"
submitFeedTx1=$($chainlinkCMD submit-feed-data feedid1 "feed 1 test data" "signatures_alice" "$alicePK" --from alice --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
submitFeedTx1Resp=$(echo "$submitFeedTx1" | jq '.height')
if [ "$submitFeedTx1Resp" == "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #1: $submitFeedTx1"
fi

# cHeCk If AlIcE gOt ThE rEwArD
echo "checking alice's reward distribution #1"
aliceCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
if [ "$aliceCurrBal" != "\"1000094\"" ]
then
  errorAndExit "Error in reward distribution for alice; expected \"1000094\", got $aliceCurrBal"
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
badSubmitFeedTx=$($chainlinkCMD submit-feed-data feedid1 "feed 1 test data" "signatures_bob" "$bobPK" --from bob --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
badSubmitFeedTxResp=$(echo "$badSubmitFeedTx" | jq '.raw_log')
if [ "$badSubmitFeedTxResp" != "\"submitter is not a valid data provider: unauthorized\"" ]
then
  errorAndExit "Error in sending bad feed data: $badSubmitFeedTx"
fi

##############

# aDd BoB aS dAtA pRoViDeR
echo "adding bob as a data provider"
addBobTx=$($chainlinkCMD add-data-provider feedid1 $bobAddr $bobPK --from alice --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
addBobTxResp=$(echo $addBobTx | jq '.height')
if [ "$addBobTxResp" == "\"0\"" ]
then
  errorAndExit "Error in adding bob as a data provider: $addBobTx"
fi

# aDd bob aS cHaInLiNk oRaClE iN aCcOuNt StOrE
echo "adding bob chainlink account"
addChainlinkAccountTx=$(chainlinkd tx chainlink add-chainlink-account "bobChainlinkPubKey" "bobChainlinkSigningKey" --from bob --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
sleep 1
addChainlinkAccountTxResp=$(echo ${addChainlinkAccountTx#*\]} | jq '.height')
if [ "$addChainlinkAccountTxResp" == "\"0\"" ]
then
  errorAndExit "Error in adding bob's chainlink account: $addChainlinkAccountTx"
fi
echo "added bob account successfully..."

# uPdAtE fEeD rEwArD
echo "updating feed reward to $newFeedReward"
newFeedReward=10
updateFeedReward=$($chainlinkCMD set-feed-reward feedid1 $newFeedReward "" --from alice --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
updateFeedRewardResp=$(echo "$updateFeedReward" | jq '.height')
if [ "$updateFeedRewardResp" == "\"0\"" ]
then
  errorAndExit "Error in updating feed reward: $updateFeedReward"
fi

# sUbMiT fEeD dAtA bY bOb
echo "submitting feed data by bob"
submitFeedTx2=$($chainlinkCMD submit-feed-data feedid1 "feed 1 test data" "signatures_bob" "$bobPK" --from bob --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
submitFeedTx2Resp=$(echo "$submitFeedTx2" | jq '.height')
if [ "$submitFeedTx2Resp" == "\"0\"" ]
then
  errorAndExit "Error in submitting feed data #2: $submitFeedTx2"
fi

# cHeCk If AlIcE gOt ThE uPdAtEd ReWaRd
echo "checking alice's reward distribution #2"
aliceCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show alice -a) --denom link --output json | jq '.amount')
if [ "$aliceCurrBal" != "\"1000088\"" ]
then
  errorAnd


  Exit "Error in reward distribution; expected \"1000088\", got $aliceCurrBal"
fi

# bOb ShOuLd NoW gEt rEwArD
echo "checking bob's reward distribution #2"
bobCurrBal=$(chainlinkd query bank balances $(chainlinkd keys show bob -a) --denom link --output json | jq '.amount')
if [ "$bobCurrBal" != "\"1000007\"" ]
then
  errorAndExit "Error in reward distribution; expected \"1000007\", got $bobCurrBal"
fi

### ~~~ BEGIN FEED EDIT TESTS ~~~ ###

# fEeD oWnEr WiLl eDiT
# eDiT sUbMiSsIoN cOuNt 
echo "edit feed's submission count by alice"
setSubmissionCountTx=$($chainlinkCMD set-submission-count feedid1 10 --from alice --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
setSubmissionCountTxResp=$(echo "$setSubmissionCountTx" | jq '.height')
if [ "$setSubmissionCountTxResp" == "\"0\"" ]
then
  errorAndExit "Error in editing submission count by alice: $setSubmissionCountTx"
fi

# eDiT hEaRtBeAt TrIgGeR
echo "edit feed's heartbeat trigger by alice"
setHeartbeatTriggerTx=$($chainlinkCMD set-heartbeat-trigger feedid1 20 --from alice --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
setHeartbeatTriggerTxResp=$(echo "$setHeartbeatTriggerTx" | jq '.height')
if [ "$setHeartbeatTriggerTxResp" == "\"0\"" ]
then
  errorAndExit "Error in editing heartbeat trigger by alice: $setHeartbeatTriggerTx"
fi

# eDiT dEvIaTiOn ThReShOlD tRiGgEr
echo "edit feed's deviation threshold trigger by alice"
setDeviationThresholdTriggerTx=$($chainlinkCMD set-deviation-threshold-trigger feedid1 30 --from alice --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
setDeviationThresholdTriggerTxResp=$(echo "$setDeviationThresholdTriggerTx" | jq '.height')
if [ "$setHeartbeatTriggerTxResp" == "\"0\"" ]
then
  errorAndExit "Error in editing deviation threshold tigger by alice: $setDeviationThresholdTriggerTx"
fi

# eDiT fEeD rEwArD
echo "edit feed's reward by alice"
setFeedRewardTx=$($chainlinkCMD set-feed-reward feedid1 40 "" --from alice --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
setFeedRewardTxResp=$(echo "$setFeedRewardTx" | jq '.height')
if [ "$setFeedRewardTxResp" == "\"0\"" ]
then
  errorAndExit "Error in editing feed reward by alice: $setFeedRewardTx"
fi

# nOn-FeEd OwMeR cAnNoT eDiT
# eDiT sUbMiSsIoN cOuNt 
echo "edit feed's submission count by cerlo"
setSubmissionCountTx=$($chainlinkCMD set-submission-count feedid1 10 --from cerlo --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
setSubmissionCountTxResp=$(echo "$setSubmissionCountTx" | jq '.height')
if [ "$setSubmissionCountTxResp" != "\"0\"" ]
then
  errorAndExit "Error in incorrect edit of submission count by cerlo: $setSubmissionCountTx"
fi

# eDiT hEaRtBeAt TrIgGeR
echo "edit feed's heartbeat trigger by cerlo"
setHeartbeatTriggerTx=$($chainlinkCMD set-heartbeat-trigger feedid1 20 --from cerlo --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
setHeartbeatTriggerTxResp=$(echo "$setHeartbeatTriggerTx" | jq '.height')
if [ "$setHeartbeatTriggerTxResp" != "\"0\"" ]
then
  errorAndExit "Error in incorrect edit of heartbeat trigger by cerlo: $setHeartbeatTriggerTx"
fi

# eDiT dEvIaTiOn ThReShOlD tRiGgEr
echo "edit feed's deviation threshold trigger by cerlo"
setDeviationThresholdTriggerTx=$($chainlinkCMD set-deviation-threshold-trigger feedid1 30 --from cerlo --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
setDeviationThresholdTriggerTxResp=$(echo "$setDeviationThresholdTriggerTx" | jq '.height')
if [ "$setDeviationThresholdTriggerTxResp" != "\"0\"" ]
then
  errorAndExit "Error in incorrect edit of deviation threshold trigger by cerlo: $setDeviationThresholdTriggerTx"
fi

# eDiT fEeD rEwArD
echo "edit feed's reward by cerlo"
setFeedRewardTx=$($chainlinkCMD set-feed-reward feedid1 40 "" --from cerlo --keyring-backend test --chain-id testchain --fees 3link <<< 'y\n')
setFeedRewardTxResp=$(echo "$setFeedRewardTx" | jq '.height')
if [ "$setFeedRewardTxResp" != "\"0\"" ]
then
  errorAndExit "Error in incorrect edit of feed reward by cerlo: $setFeedRewardTx"
fi


pkill chainlinkd
echo "Chainlink module FEED test has exited successfully."
exit 0
