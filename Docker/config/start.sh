#!/bin/bash
# Name: start.sh
# Purpose: Run multiple sawtooth nodes on multiple machines
# ----------------------------------------------------

DOCKER_COMPOSE_PATH=/home/jfp/SignerNode_Thesis/Docker/config
BOOTSTRAP=51.83.75.29
SERVERS=(51.83.75.29 51.83.75.29 51.83.75.29 51.83.75.29 51.83.75.29)

ID=1
API_PORT=7000
PORT=46000
CRYPTO_PORT=10000
SMARTCONTRACT_PORT=11000

#
#Start up signer node peer nodes
#------------------------------------------------------------

 CMD="bash -c '
    docker-compose -p ${ID} -f ${DOCKER_COMPOSE_PATH}/bootstrap.yaml up --detach
    '
  "
  echo $CMD | ssh -t jfp@${BOOTSTRAP} bash

sleep 1

#
#Start up signer node peer nodes
#------------------------------------------------------------
for s in ${SERVERS[@]}
do
  CMD="bash -c '
    export ID=${ID} &&
    export BOOTSTRAP=${BOOTSTRAP} &&
    export IP=${s} &&
    export PORT=${PORT} &&
    export API_PORT=${API_PORT} &&
    export CRYPTO_PORT=${CRYPTO_PORT} &&
    export SMARTCONTRACT_PORT=${SMARTCONTRACT_PORT} &&
    docker-compose -p ${ID} -f ${DOCKER_COMPOSE_PATH}/peer.yaml up --detach
    '
  "
  (echo $CMD | ssh -t jfp@${s} bash)
  echo "- \"${s}:${API_PORT}\"" >> peers.txt

  ID=$((ID + 1))
  PORT=$((PORT + 1))
  API_PORT=$((API_PORT + 1))
  CRYPTO_PORT=$((CRYPTO_PORT + 1))
  SMARTCONTRACT_PORT=$((SMARTCONTRACT_PORT + 1))
done