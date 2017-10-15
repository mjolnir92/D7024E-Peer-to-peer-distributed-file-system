#!/bin/bash
node_prefix="kademlia"
num_nodes=3
num_nodes=$(expr $num_nodes - 1)

echo "Starting nodes"
container_name=()
# set up the first node alone
i=0
container_name[i]=$(docker run -d kademlia)
container_ip[i]=$(docker inspect ${container_name[$i]} | jq -r '.[0].NetworkSettings.Networks.bridge.IPAddress')
for i in $(seq 1 $num_nodes); do
	container_name[i]=$(docker run -d kademlia -j ${container_ip[$(expr $i - 1)]}:1200)
	container_ip[i]=$(docker inspect ${container_name[$i]} | jq -r '.[0].NetworkSettings.Networks.bridge.IPAddress')
done
echo "Done"
