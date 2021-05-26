#!/bin/bash
# Name: start.sh
# Purpose: Run multiple sawtooth nodes on multiple machines
# ----------------------------------------------------

DOCKER_COMPOSE_PATH=/home/ubuntu/SignerNode_Thesis/Docker/config
BOOTSTRAP=146.59.230.10
#SERVERS=(146.59.230.10 146.59.229.32 146.59.228.136 146.59.228.136 141.95.18.233 135.125.207.187 141.95.18.234 141.95.18.234 51.79.53.77 51.79.53.77 51.79.53.170 51.79.53.170 139.99.171.133 139.99.171.133 139.99.236.2 139.99.236.2)

#Complete set of servers for 16 nodes
#SERVERS=(146.59.230.10 146.59.229.32 146.59.228.136 141.95.18.233 135.125.207.187 141.95.18.234 51.79.53.77 51.79.53.77 51.79.53.170 51.79.53.170 51.79.53.170 139.99.171.133 139.99.171.133 139.99.236.2 139.99.236.2 139.99.236.2)

#Crash Fault 1 (France)
#SERVERS=(146.59.230.10 146.59.228.136 141.95.18.233 135.125.207.187 141.95.18.234 51.79.53.77 51.79.53.77 51.79.53.170 51.79.53.170 51.79.53.170 139.99.171.133 139.99.171.133 139.99.236.2 139.99.236.2 139.99.236.2)

#Crash Fault 2 (France, Australia)
#SERVERS=(146.59.230.10 146.59.228.136 141.95.18.233 135.125.207.187 141.95.18.234 51.79.53.77 51.79.53.77 51.79.53.170 51.79.53.170 51.79.53.170 139.99.171.133 139.99.171.133 139.99.236.2 139.99.236.2)

#Crash Fault 3 (France, Australia, Canada)
#SERVERS=(146.59.230.10 146.59.228.136 141.95.18.233 135.125.207.187 141.95.18.234 51.79.53.77 51.79.53.77 51.79.53.170 51.79.53.170 139.99.171.133 139.99.171.133 139.99.236.2 139.99.236.2)

#Crash Fault 4 (France,Germany, Australia, Canada)
#SERVERS=(146.59.230.10 146.59.228.136 141.95.18.233 135.125.207.187  51.79.53.77 51.79.53.77 51.79.53.170 51.79.53.170 139.99.171.133 139.99.171.133 139.99.236.2 139.99.236.2)

#Crash Fault 5 (France,Germany, Australia (2), Canada)
SERVERS=(146.59.230.10 146.59.228.136 141.95.18.233 135.125.207.187 51.79.53.77 51.79.53.77 51.79.53.170 51.79.53.170 139.99.171.133 139.99.171.133 139.99.236.2)

BYZANTINE=(146.59.229.32 141.95.18.234 51.79.53.77 51.79.53.170 139.99.236.2)

ID=1
API_PORT=7000
PORT=46000
CRYPTO_PORT=10000
SMARTCONTRACT_PORT=11000

N=16
T=11
SCHEME=TRSA2048Pessimistic

#
#Start up signer node peer nodes
#------------------------------------------------------------

 CMD="bash -c '
    docker-compose -p ${ID} -f ${DOCKER_COMPOSE_PATH}/bootstrap.yaml up --detach
    '
  "
  echo $CMD | ssh -t ubuntu@${BOOTSTRAP} bash

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
    export T=${T} &&
    export N=${N} &&
    export SCHEME=${SCHEME} &&
    docker-compose -p ${ID} -f ${DOCKER_COMPOSE_PATH}/peer.yaml up --detach
    '
  "
  (echo $CMD | ssh -t ubuntu@${s} bash)
  echo "- \"${s}:${API_PORT}\"" >> peers.txt

  ID=$((ID + 1))
  PORT=$((PORT + 1))
  API_PORT=$((API_PORT + 1))
  CRYPTO_PORT=$((CRYPTO_PORT + 1))
  SMARTCONTRACT_PORT=$((SMARTCONTRACT_PORT + 1))
done


#
#Start up Byzantine nodes
#------------------------------------------------------------
for s in ${BYZANTINE[@]}
do
  CMD="bash -c '
    export ID=${ID} &&
    export BOOTSTRAP=${BOOTSTRAP} &&
    export IP=${s} &&
    export PORT=${PORT} &&
    export API_PORT=${API_PORT} &&
    export CRYPTO_PORT=${CRYPTO_PORT} &&
    export SMARTCONTRACT_PORT=${SMARTCONTRACT_PORT} &&
    export T=${T} &&
    export N=${N} &&
    export SCHEME=${SCHEME} &&
    docker-compose -p ${ID} -f ${DOCKER_COMPOSE_PATH}/byz.yaml up --detach
    '
  "
  (echo $CMD | ssh -t ubuntu@${s} bash)

  ID=$((ID + 1))
  PORT=$((PORT + 1))
  API_PORT=$((API_PORT + 1))
  CRYPTO_PORT=$((CRYPTO_PORT + 1))
  SMARTCONTRACT_PORT=$((SMARTCONTRACT_PORT + 1))
done