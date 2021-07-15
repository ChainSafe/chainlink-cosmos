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
goodTx1=$(chainlinkd tx chainlink addFeed feedid1 $aliceAddr 1 2 3 4 $aliceAddr,$alicePK --from alice --keyring-backend test --chain-id testchain <<< 'y\n')
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
echo $goodTx2
# "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"SubmitFeedData\"},{\"key\":\"sender\",\"value\":\"cosmos1kjl7py8a6mxlg5ase9sutn0qdg36teqpfldhhy\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos1jftpxmcunts62sf0x3d4qtk2khflkx436rf0tk\"},{\"key\":\"sender\",\"value\":\"cosmos1kjl7py8a6mxlg5ase9sutn0qdg36teqpfldhhy\"},{\"key\":\"amount\",\"value\":\"4link\"}]}]}]"
if [ "$goodTx2Resp" != "\"[{\\\"events\\\":[{\\\"type\":\\\"message\\\",\\\"attributes\\\":[{\\\"key\\\":\\\"action\\\",\\\"value\\\":\\\"SubmitFeedData\\\"},{\\\"key\\\":\\\"sender\\\",\\\"value\\\":\\\"cosmos1kjl7py8a6mxlg5ase9sutn0qdg36teqpfldhhy\\\"}]},{\\\"type\\\":\\\"transfer\\\",\\\"attributes\\\":[{\\\"key\\\":\\\"recipient\\\",\\\"value\\\":\\\""$aliceAddr"\\\"},{\\\"key\\\":\\\"sender\\\",\\\"value\\\":\\\"cosmos1kjl7py8a6mxlg5ase9sutn0qdg36teqpfldhhy\\\"},{\\\"key\\\":\\\"amount\\\",\\\"value\\\":\\\"4link\\\"}]}]}]\"" ]

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



# # start.sh
# chainlinkd tx chainlink addFeed feedid1 "$(chainlinkd keys show alice -a)" 1 2 3 4 "$(chainlinkd keys show alice -a)","$(chainlinkd keys show alice -p)" --from alice --keyring-backend test --chain-id testchain <<< 'y\n'

# chainlinkd tx chainlink submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain <<< 'y\n'



# chainlinkd tx chainlink addModuleOwner "$(chainlinkd keys show bob -a)" "$(chainlinkd keys show bob -p)" --from alice --keyring-backend test --chain-id testchain <<< 'y\n'

# chainlinkd tx chainlink addFeed feedid1 "$(chainlinkd keys show bob -a)" 1 2 3 4 "$(chainlinkd keys show bob -a),$(chainlinkd keys show bob -p)" --from alice --keyring-backend test --chain-id testchain <<< 'y\n'

# chainlinkd tx chainlink submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from bob --keyring-backend test --chain-id testchain <<< 'y\n'

# # query balance
# chainlinkd query bank balances $(chainlinkd keys show alice -a)
# chainlinkd query bank balances $(chainlinkd keys show bob -a)


# chainlinkd query tx 426EA8DD65D273BBACA4AF21C7D4220F57930752F563A2788A1D464EF611F8AA --chain-id testchain -o json
# root@fcb0a52a4afd:/chainlink# chainlinkd query tx 426EA8DD65D273BBACA4AF21C7D4220F57930752F563A2788A1D464EF611F8AA --chain-id testchain -o json
# {"height":"96","txhash":"426EA8DD65D273BBACA4AF21C7D4220F57930752F563A2788A1D464EF611F8AA","codespace":"undefined","code":111222,"data":"","raw_log":"panic message redacted to hide potentially sensitive system info: panic","logs":[],"info":"","gas_wanted":"200000","gas_used":"52301","tx":{"@type":"/cosmos.tx.v1beta1.Tx","body":{"messages":[{"@type":"/chainlink.v1beta.MsgFeedData","feedId":"feedid1","submitter":"cosmos1em2sjgsssgtch4g7av7q7ns7gcvvgr3wvxftlv","feedData":"ZmVlZCAxIHRlc3QgZGF0YQ==","signatures":["ZHVtbXk=","c2lnbmF0dXJlcw=="]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Am9eVZMDGJIlX3b7GsCzc2sQDakUs8YWRkKb5006BAhM"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"5"}],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""}},"signatures":["SzZkND4ztVJF1qtp6aerLoU7k1eIuk0iOgRyat9KGmBbr1E63yxrjfxnbGPg0FQZrFYKBb8fb3XKyv25csVTMw=="]},"timestamp":"2021-07-14T05:43:15Z"}
