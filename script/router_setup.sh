#! /bin/bash
# 参考にしたlink https://cipepser.hatenablog.com/entry/2018/06/09/004657
# ------------------------------------------------------------------------------------------#
#                                                                                           #
#   --------                  --------                  ----------                --------  #   
#   | HOST | veth0-------veth1| next |veth0--------veth1| router |veth0------veth1| host |  #
#   -------- .254          .1 -------- .254          .1 ---------- .254        .1 --------  #
#             192.168.2.0/24            192.168.1.0/24              192.168.0.0/24          #
# ------------------------------------------------------------------------------------------#
# 名前空間を区切る
ip netns add host
ip netns add router
ip netns add next

# 区切った名前空間をvethで繋ぐ
ip link add host_veth1 type veth peer name router_veth0
ip link add router_veth1 type veth peer name next_veth0
ip link add next_veth1 type veth peer name linux_veth0

# vethを各名前空間に所属させる
ip link set host_veth1 netns host
ip link set router_veth0 netns router
ip link set router_veth1 netns router
ip link set next_veth0 netns next
ip link set next_veth1 netns next

# アドレスを割り当てる
ip netns exec host ip addr add 192.168.0.1/24 dev host_veth1
ip netns exec router ip addr add 192.168.0.254/24 dev router_veth0
ip netns exec router ip addr add 192.168.1.1/24 dev router_veth1
ip netns exec next ip addr add 192.168.1.254/24 dev next_veth0
ip netns exec next ip addr add 192.168.2.1/24 dev next_veth1
ip addr add 192.168.2.254/24 dev linux_veth0

# 立ち上げ
ip netns exec host ip link set lo up
ip netns exec router ip link set lo up
ip netns exec next ip link set lo up
ip netns exec host ip link set host_veth1 up
ip netns exec router ip link set router_veth0 up
ip netns exec router ip link set router_veth1 up
ip netns exec next ip link set next_veth0 up
ip netns exec next ip link set next_veth1 up
ip link set linux_veth0 up

# ルーティングの設定
ip netns exec host ip route add default via 192.168.0.254
ip netns exec router ip route add default via 192.168.1.254
ip netns exec next ip route add default via 192.168.2.254
ip netns exec next ip route add 192.168.0.0/24 via 192.168.1.1
# ip netns exec next ip route add 192.168.1.0/24 via 192.168.1.1
ip route add 192.168.0.0/24 via 192.168.2.1
ip route add 192.168.1.0/24 via 192.168.2.1

# ip_forwardの設定
mkdir -p /etc/netns/router
cp /etc/sysctl.conf /etc/netns/router/
sysctl -w net.ipv4.ip_forward=1

