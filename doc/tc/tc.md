## Traffict Control Usage
---

#### Show current settings

tc qdisc show dev eth0

---

#### Del all settings

tc qdisc del dev eth0 root

---

#### Use netem: delay, loss, corrupt, duplicate

tc qdisc add dev eth0 root netem delay 100ms 3ms

tc qdisc add dev eth0 root netem delay 100ms 10ms distribution normal

tc qdisc add dev eth0 root netem loss 5%

tc qdisc change dev eth0 root netem corrupt 5% duplicate 1%

---

#### Use tbf

tc qdisc add dev eth0 root tbf rate 1mbit burst 32kbit latency 400ms

> * tbf: use the token buffer filter to manipulate traffic rates
> * rate: sustained maximum rate
> * burst: maximum allowed burst
> * latency: packets with higher latency get dropped

tc qdisc add dev eth0 root tbf rate 1mbit burst 10kb latency 70ms peakrate 2mbit minburst 1540

---

#### References:

[1] https://netbeez.net/blog/how-to-use-the-linux-traffic-control/
[2] https://www.cyberciti.biz/faq/linux-traffic-shaping-using-tc-to-control-http-traffic/
[3] http://lartc.org/howto/lartc.qdisc.html
[4] http://manpages.ubuntu.com/manpages/xenial/man8/tc.8.html
[5] https://unix.stackexchange.com/questions/100785/bucket-size-in-tbf
