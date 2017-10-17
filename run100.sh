#!/bin/bash
node_prefix="kademlia"
num_nodes=100
num_nodes=$(expr $num_nodes - 1)

container_name=()
# set up the first node alone
i=0
echo "Starting node $i"
container_name[i]=$(docker run -d kademlia)
container_ip[i]=$(docker inspect ${container_name[$i]} | jq -r '.[0].NetworkSettings.Networks.bridge.IPAddress')
for i in $(seq 1 $num_nodes); do
	echo "Starting node $i"
	container_name[i]=$(docker run -d kademlia -j ${container_ip[0]}:1200)
	container_ip[i]=$(docker inspect ${container_name[$i]} | jq -r '.[0].NetworkSettings.Networks.bridge.IPAddress')
done
echo "Done"
