#!/bin/bash
node_prefix="kademlia"
num_nodes=20
num_nodes=$(expr $num_nodes - 1)

container_name=()
# set up the first node alone
i=0
echo "Starting node $i"
container_name[i]=$(docker run -d kademlia)
container_ip[i]=$(docker inspect ${container_name[$i]} | jq -r '.[0].NetworkSettings.Networks.bridge.IPAddress')
for i in $(seq 1 $num_nodes); do
	sleep 0.2
	echo "Starting node $i"
	container_name[i]=$(docker run -d kademlia -j ${container_ip[0]}:1200)
	container_ip[i]=$(docker inspect ${container_name[$i]} | jq -r '.[0].NetworkSettings.Networks.bridge.IPAddress')
done

echo "store"
key=$(./kdfs -s ${container_ip[1]}:8080 store ./README.md)
echo "cat"
cat=$(./kdfs -s ${container_ip[2]}:8080 cat $key)
if [[ $cat ]]; then
	echo "ok (found data)"
	echo $cat
else
	echo "fail (no data found)"
fi
echo "pin"
./kdfs -s ${container_ip[3]}:8080 pin $key
echo "sleep"
sleep 30
echo "cat"
cat=$(./kdfs -s ${container_ip[4]}:8080 cat $key)
if [[ $cat ]]; then
	echo "ok (data stayed after sleeping)"
	echo $cat
else
	echo "fail (data didn't stay after pinning)"
fi
echo "unpin"
./kdfs -s ${container_ip[5]}:8080 unpin $key
echo "sleep"
sleep 30
echo "cat"
cat=$(./kdfs -s ${container_ip[6]}:8080 cat $key)
if [[ $cat ]]; then
	echo "fail (data stayed after unpinning)"
	echo $cat
else
	echo "ok (data deleted correctly)"
fi
