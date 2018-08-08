tc qdisc del dev eth0 root

#tc qdisc add dev eth0 root handle 1: htb default 10
#tc class add dev eth0 parent 1: classid 1:30 htb rate 10mbps
#tc class add dev eth0 parent 1: classid 1:1 htb rate 10mbps
#tc class add dev eth0 parent 1: classid 1:2 htb rate 10mbps
#tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst 172.17.0.3/0 flowid 1:1
#tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip src 172.17.0.3/0 flowid 1:2

tc qdisc add dev eth0 root handle 1: htb default 10

tc class add dev eth0 parent 1: classid 1:1 htb rate 10mbps
tc class add dev eth0 parent 1:1 classid 1:10 htb rate 10mbps
tc class add dev eth0 parent 1:1 classid 1:20 htb rate 10mbps
tc class add dev eth0 parent 1:1 classid 1:30 htb rate 10mbps
tc class add dev eth0 parent 1:1 classid 1:40 htb rate 10mbps

tc qdisc add dev eth0 parent 1:10 handle 10: netem delay 100ms 5ms
tc qdisc add dev eth0 parent 1:20 handle 20: tbf rate 1mbit burst 32kbit latency 100ms
tc qdisc add dev eth0 parent 1:30 handle 30: tbf rate 2mbit burst 32kbit latency 100ms
tc qdisc add dev eth0 parent 1:40 handle 40: sfq

tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip src 172.17.0.3/32 flowid 1:20
tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst 172.17.0.3/32 flowid 1:20

tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip src 172.17.0.4/32 flowid 1:30
tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst 172.17.0.4/32 flowid 1:30
