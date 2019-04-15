#!/usr/bin/env bash

function printUsage() {
    echo "Usage: "
    echo "  network.sh <action>"
    echo "    <action>: 'up', 'down' or 'generate'"
    echo "      - 'up': bring up the network with docker-compose up"
    echo "      - 'down': bring down the network with docker-compose down"
    echo "      - 'generate [crypto|channel]': generate required certificates and the genesis block"
}

CHANNEL_NAME="docshare"
ACTION="${1}"
shift

function generateCerts() {
    which ./bin/cryptogen
    if [ "$?" -ne 0 ]; then
        echo "cryptogen is absent. exiting..."
        echo "to resolve set the PATH variable"
        exit 1
    fi
    echo "#########  Generating Certificates ##############"
    if [ -d "crypto-config" ]; then
        rm -rf crypto-config
    fi
    mkdir -p crypto-config
    ./bin/cryptogen generate --config=./crypto-config.yaml
    if [ "$?" -ne 0 ]; then
        echo "failed to generate certificates..."
        exit 1
    fi
    echo
}

function generateChannelArtifacts() {
    which ./bin/configtxgen
    if [ "$?" -ne 0 ]; then
        echo "configtxgen tool not found. exiting..."
        echo "to resolve set the PATH variable"
        exit 1
    fi
    echo "#########  Generating Orderer Genesis block ##############"
    mkdir -p artifacts
    FABRIC_CFG_PATH=$PWD ./bin/configtxgen -profile DocShare -outputBlock ./artifacts/orderer.genesis.block
    if [ "$?" -ne 0 ]; then
        echo "failed to generate orderer genesis block..."
        exit 1
    fi
    echo
    echo "### Generating channel configuration transaction ###"
    FABRIC_CFG_PATH=$PWD ./bin/configtxgen -profile DocShare -outputCreateChannelTx ./artifacts/docshare.channel.tx -channelID $CHANNEL_NAME
    if [ "$?" -ne 0 ]; then
        echo "failed to generate channel configuration transaction..."
        exit 1
    fi
    echo
    echo "### Generating anchor peer transaction ###"
    FABRIC_CFG_PATH=$PWD ./bin/configtxgen -profile DocShare -outputAnchorPeersUpdate ./artifacts/org1.docshare.anchors.tx -channelID $CHANNEL_NAME -asOrg DocShareOrganization1
    if [ "$?" -ne 0 ]; then
        echo "failed to generate channel configuration transaction..."
        exit 1
    fi
}

function networkUp() {
    # generate artifacts if they don't exist
    if [ ! -d "crypto-config" ]; then
        generateCerts
        generateChannelArtifacts
    fi
    docker-compose -f docker-compose.yaml up -d
    if [ $? -ne 0 ]; then
        echo "ERROR !!!! Unable to start network"
        exit 1
    fi
}

function networkDown() {
    docker-compose -f docker-compose.yaml stop
    if [ $? -ne 0 ]; then
        echo "ERROR !!!! Test failed"
        exit 1
    fi
}

if [ "${ACTION}" == "up" ]; then
    networkUp
elif [ "${ACTION}" == "down" ]; then
    networkDown
elif [ "${ACTION}" == "generate" ]; then
    if [ "${1}" == ""  ]; then
        generateCerts
        generateChannelArtifacts
    elif [ "${1}" == "crypto" ]; then
        generateCerts
    elif [ "${1}" == "channel" ]; then
        generateChannelArtifacts
    else
        printUsage
        exit 1
    fi
else
    printUsage
    exit 1
fi
