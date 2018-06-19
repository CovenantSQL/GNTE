# 网络拓扑描述语言
build a cluster with traffic control
## YAML形式
```
group:
  -
    name: china
    nodes:
        - 192.168.1.1/24
        - 10.1.1.1/20
        - 11.1.1.1/20
    delay: "100ms 10ms 30%"
    loss: "1% 10%"
  -
    name: eu
    nodes:
        - 12.1.1.1/20
        - 13.1.1.1/20
        - 14.1.1.1/20
    delay: "10ms 5ms 30%"
    loss: "1% 10%"
  -
    name: jpn
    nodes:
        - 15.1.1.1/20
    delay: "100ms 10ms 30%"
    duplicate: "1%"
    rate: "100mbit"

network:
  -
    groups:
        - china
        - eu
    delay: "200ms 10ms 1%"
    corrupt: "0.2%"
    rate: "10mbit"

  -
    groups:
        - china
        - jpn
    delay: "100ms 10ms 1%"
    rate: "10mbit"

  -
    groups:
        - jpn
        - eu
    delay: "30ms 5ms 1%"
    rate: "100mbit"
```

## 关键词说明
### node
描述单机的基础属性，可以不写，表示机器无限制。规则只应用于出口(上行)

### group
描述单机组成的集群或区域，可以嵌套group

### network
整个网络的描述入口，写入这里的会最终翻译成TC规则

### ip
后面跟着机器IP，可以用CIDR描述多个机器

### delay
delay 100ms ,同时,大约有 30% 的包会延迟 ± 10ms 发送

### loss
丢包率，同TC参数

### duplicate
重复数据包概率，同TC参数

### rate
限速度参数，根据TC的定义，这块的限速不是精准限速，会在这个值浮动一些

## 规则冲突处理
同一个级别的网络，如node之间，group之间，以文件顺序最后写的为准。

不同级别的网络为嵌套规则，叠加生效
