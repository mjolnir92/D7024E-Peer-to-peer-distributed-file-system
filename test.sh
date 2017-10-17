#!/bin/bash
echo "store"
key=$(./kdfs -s ${CONTAINER_IP[1]}:8080 store ./README.md)
echo "cat"
cat=$(./kdfs -s ${CONTAINER_IP[98]}:8080 cat $key)
if [[ $cat ]]; then
	echo "ok (found data)"
	echo $cat
else
	echo "fail (no data found)"
fi
echo "pin"
./kdfs -s ${CONTAINER_IP[15]}:8080 pin $key
echo "sleep"
sleep 30
echo "cat"
cat=$(./kdfs -s ${CONTAINER_IP[97]}:8080 cat $key)
if [[ $cat ]]; then
	echo "ok (data stayed after sleeping)"
	echo $cat
else
	echo "fail (data didn't stay after pinning)"
fi
echo "unpin"
./kdfs -s ${CONTAINER_IP[25]}:8080 unpin $key
echo "sleep"
sleep 30
echo "cat"
cat=$(./kdfs -s ${CONTAINER_IP[96]}:8080 cat $key)
if [[ $cat ]]; then
	echo "fail (data stayed after unpinning)"
	echo $cat
else
	echo "ok (data deleted correctly)"
fi
