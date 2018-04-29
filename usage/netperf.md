## Netperf Usage
---
*Run First*

Server

netserver -4

---

netperf -H 172.17.0.2 -l 10 -f M

-H host
-l test duration (>0 secs) (<0 bytes|trans)
-f output units: 'M' means 10^6 Bytes

---

netperf -H 172.17.0.2 -l 5 -f g -- -m 1024 -M 1024 -s 1024

-f output units: 'g' means 10^9 bps
-m send size for both side
-M recv size for both side
-s send & recv socket buff size for *local* side
-S send & recv socket buff size for *remote* side
