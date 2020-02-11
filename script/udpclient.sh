#! /bin/bash

ip netns add server
ip netns add client

ip link add server_veth0 type veth peer client_veth0

ip link set server_veth0 netns server
ip link set client_veth0 netns client

ip netns exec server ip addr add 192.168.0.2/24 dev server_veth0
ip netns exec client ip addr add 192.168.0.3/24 dev client_veth0

ip netns exec server ip link set lo up
ip netns exec client ip link set lo up
ip netns exec server ip link set server_veth0 up
ip netns exec client ip link set client_veth0 up

