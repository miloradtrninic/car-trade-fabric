#!/bin/bash


echo
echo " ____    _____      _      ____    _____ "
echo "/ ___|  |_   _|    / \    |  _ \  |_   _|"
echo "\___ \    | |     / _ \   | |_) |   | |  "
echo " ___) |   | |    / ___ \  |  _ <    | |  "
echo "|____/    |_|   /_/   \_\ |_| \_\   |_|  "
echo
echo "Build your first network (BYFN) end-to-end test"
echo
CHANNEL_NAME="$1"
DELAY="$2"
LANGUAGE="$3"
TIMEOUT="$4"
VERBOSE="$5"
: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="10"}
: ${LANGUAGE:="golang"}
: ${TIMEOUT:="10"}
: ${VERBOSE:="false"}
LANGUAGE=`echo "$LANGUAGE" | tr [:upper:] [:lower:]`
COUNTER=1
MAX_RETRY=10

CC_SRC_PATH="github.com/chaincode/car-sales/go/"
if [ "$LANGUAGE" = "node" ]; then
	CC_SRC_PATH="/opt/gopath/src/github.com/chaincode/car-sales/node/"
fi

if [ "$LANGUAGE" = "java" ]; then
	CC_SRC_PATH="/opt/gopath/src/github.com/chaincode/car-sales/java/"
fi

echo "Channel name : "$CHANNEL_NAME

# import utils
. scripts/utils.sh

createChannel() {
	setGlobals 0 1

	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
                set -x
		peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx >&log.txt
		res=$?
                set +x
	else
				set -x
		peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
		res=$?
				set +x
	fi
	cat log.txt
	verifyResult $res "Channel creation failed"
	echo "===================== Channel '$CHANNEL_NAME' created ===================== "
	echo
}

joinChannel () {
	for org in 1 2; do
	    for peer in 0 1 2; do
			joinChannelWithRetry $peer $org
			echo "===================== peer${peer}.org${org} joined channel '$CHANNEL_NAME' ===================== "
			sleep $DELAY
			echo
	    done
	done
	for org in 3 4; do
	    for peer in 0 1 2 3; do
			joinChannelWithRetry $peer $org
			echo "===================== peer${peer}.org${org} joined channel '$CHANNEL_NAME' ===================== "
			sleep $DELAY
			echo
	    done
	done
}

## Create channel
echo "Creating channel..."
createChannel

## Join all the peers to the channel
echo "Having all peers join the channel..."
joinChannel

## Set the anchor peers for each org in the channel
echo "Updating anchor peers for org1..."
updateAnchorPeers 0 1
echo "Updating anchor peers for org2..."
updateAnchorPeers 0 2
echo "Updating anchor peers for org3..."
updateAnchorPeers 0 3
echo "Updating anchor peers for org4..."
updateAnchorPeers 0 4

## Install chaincode 
echo "Installing chaincode on peer0.org1..."
installChaincode 0 1
echo "Installing chaincode on peer0.org1..."
installChaincode 1 1
echo "Installing chaincode on peer0.org1..."
installChaincode 2 1

echo "Installing chaincode on peer0.org2..."
installChaincode 0 2
echo "Installing chaincode on peer1.org2..."
installChaincode 1 2
echo "Installing chaincode on peer2.org2..."
installChaincode 2 2

echo "Installing chaincode on peer0.org3..."
installChaincode 0 3
echo "Installing chaincode on peer1.org3..."
installChaincode 1 3
echo "Installing chaincode on peer2.org3..."
installChaincode 2 3
echo "Installing chaincode on peer3.org3..."
installChaincode 3 3

echo "Installing chaincode on peer0.org4..."
installChaincode 0 4
echo "Installing chaincode on peer1.org4..."
installChaincode 1 4
echo "Installing chaincode on peer2.org4..."
installChaincode 2 4
echo "Installing chaincode on peer3.org4..."
installChaincode 3 4

export WAIT_INSTANTIATE=60
echo "Fabric INSTANTIATE timeout ${WAIT_INSTANTIATE} "
sleep ${WAIT_INSTANTIATE}


# Instantiate chaincode
echo "Instantiating chaincode on peer0.org1..."
instantiateChaincode 0 1


export WAIT_INVOKE=60
echo "Fabric invoke timeout ${WAIT_INSTANTIATE} "
sleep ${WAIT_INVOKE}

# Query chaincode on peer0.org1
#echo "Querying chaincode on peer0.org1..."
#chaincodeQuery 0 1 100

# Invoke chaincode on peer0.org1 and peer0.org2
echo "Sending invoke transaction on peer0.org1"
chaincodeInvoke 0 1



echo
echo "========= All GOOD, BYFN execution completed =========== "
echo

echo
echo " _____   _   _   ____   "
echo "| ____| | \ | | |  _ \  "
echo "|  _|   |  \| | | | | | "
echo "| |___  | |\  | | |_| | "
echo "|_____| |_| \_| |____/  "
echo

exit 0
