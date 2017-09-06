#!/usr/bin/env bash

# Debug
# set -x

# Function to clean up all the containers
function cleanup_and_exit {
  echo "[+] Cleaning up containers"
  echo "[+] Stopping containers"
  docker stop iyomongo-cfg0 iyomongo-cfg1 iyomongo-cfg2
  docker stop iyomongo-ps0 iyomongo-ps1 iyomongo-ps2
  docker stop iyomongo-ss00 iyomongo-ss01 iyomongo-ss02
  docker stop iyomongo-mongos
  echo "[+] Containers stopped"
  echo "[+] Removing containers"
  docker rm iyomongo-cfg0 iyomongo-cfg1 iyomongo-cfg2
  docker rm iyomongo-ps0 iyomongo-ps1 iyomongo-ps2
  docker rm iyomongo-ss00 iyomongo-ss01 iyomongo-ss02
  docker rm iyomongo-mongos
  echo "[+] Containers removed"
  echo "[+] Exiting"
  exit
}

# Clean up if we encounter an error
trap cleanup_and_exit ERR

# If the script is run with 'clean' as first argument, just remove the dockers
if [ "$1" == "clean" ]; then cleanup_and_exit; fi

# Create the config set
# By default mongo servers with the 'configsvr' flag set run on port 27019. Change
# it here to the default port so we don't have to remember when connecting a shell
# later
echo "[+] Creating config set servers"
docker run -d --name iyomongo-cfg0 mongo:3.4 --replSet cfg --configsvr --port 27017
docker run -d --name iyomongo-cfg1 mongo:3.4 --replSet cfg --configsvr --port 27017
docker run -d --name iyomongo-cfg2 mongo:3.4 --replSet cfg --configsvr --port 27017
echo "[+] Config set members created"

echo "[+] Getting config set members ip addresses"
ip_cfg0=$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" iyomongo-cfg0)
ip_cfg1=$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" iyomongo-cfg1)
ip_cfg2=$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" iyomongo-cfg2)

echo "[+] Initializing config set"
# Make sure the config members are fully initialized before we connect
sleep 3
docker exec iyomongo-cfg0 mongo --quiet --eval "rs.initiate({ _id : 'cfg', members : [ {_id: 0, host: '$ip_cfg0:27017'}, {_id: 1, host: '$ip_cfg1:27017'}, {_id: 2, host: '$ip_cfg2:27017'} ] });" > /dev/null
echo "[+] Config set initialized"

# Create the primary shard set
# By default mongo servers with the 'shardsvr' flag set run on port 27018. Change
# it here to the default port so we don't have to remember when connecting a shell
# later
echo "[+] Creating primary shard servers"
docker run -d --name iyomongo-ps0 mongo:3.4 --replSet ps --shardsvr --port 27017
docker run -d --name iyomongo-ps1 mongo:3.4 --replSet ps --shardsvr --port 27017
docker run -d --name iyomongo-ps2 mongo:3.4 --replSet ps --shardsvr --port 27017
echo "[+] Primary shard members created"

echo "[+] Getting primary shard members ip addresses"
ip_ps0=$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" iyomongo-ps0)
ip_ps1=$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" iyomongo-ps1)
ip_ps2=$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" iyomongo-ps2)

echo "[+] Initalizing primary shard"
# Make sure the primay shard members are fully initialized before we connect
sleep 3
docker exec iyomongo-ps0 mongo --quiet --eval "rs.initiate({ _id : 'ps', members : [ {_id: 0, host: '$ip_ps0:27017'}, {_id: 1, host: '$ip_ps1:27017'}, {_id: 2, host: '$ip_ps2:27017'} ] });" > /dev/null
echo "[+] Primary shard initialized"

# Create the secondary shard set
# Like the primary shard, make sure they are started on port 27017
echo "[+] Creating secondary shard servers"
docker run -d --name iyomongo-ss00 mongo:3.4 --replSet ss0 --shardsvr --port 27017
docker run -d --name iyomongo-ss01 mongo:3.4 --replSet ss0 --shardsvr --port 27017
docker run -d --name iyomongo-ss02 mongo:3.4 --replSet ss0 --shardsvr --port 27017
echo "[+] Secondary shard members created"

echo "[+] Getting secondary shard members ip addresses"
ip_ss00=$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" iyomongo-ss00)
ip_ss01=$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" iyomongo-ss01)
ip_ss02=$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" iyomongo-ss02)

echo "[+] Initializing secondary shard"
# Make sure the secondary shard members are fully initialized before we connect
sleep 3
docker exec iyomongo-ss00 mongo --quiet --eval "rs.initiate({ _id : 'ss0', members : [ {_id: 0, host: '$ip_ss00:27017'}, {_id: 1, host: '$ip_ss01:27017'}, {_id: 2, host: '$ip_ss02:27017'} ] });" > /dev/null
echo "[+] Secondary shard initialized"

# Init mongos
# Make sure to forward this port so we can connect to localhost from applications
echo "[+] Initializing mongos"
docker run -d --name iyomongo-mongos -p 27017:27017 mongo:3.4 mongos --configdb "cfg/$ip_cfg0:27017,$ip_cfg1:27017,$ip_cfg2:27017"
echo "[+] Giving mongos some time to fully initialize"
sleep 15
echo "[+] Mongos initialized"

# Add shards to mongos
echo "[+] Adding primary shard to mongos"
docker exec iyomongo-mongos mongo --quiet --eval "sh.addShard(\"ps/$ip_ps0:27017,$ip_ps1:27017,$ip_ps2:27017\")"
echo "[+] Primary shard added"

echo "[+] Adding secondary shard to mongos"
docker exec iyomongo-mongos mongo --quiet --eval "sh.addShard(\"ss0/$ip_ss00:27017,$ip_ss01:27017,$ip_ss02:27017\")"
echo "[+] Secondary shard added"

echo "[+] Database setup ready to receive data"
echo "[+] Sharding must still be enabled"
