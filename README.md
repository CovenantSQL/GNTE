<img src="logo/logo.jpeg" width=200>

# GNTE
GNTE(Global Network Topology Emulator) is a docker based all-in-one emulator which emulating unstable global network, such as random delay, packet loss, etc.

## Before Use
install docker

## Build and Run
### 1. build docker image
Clone this repo and run ```build.sh```, there should be an image named ```ns``` in your docker environment.

### 2. modify network definition file
Edit ```example.yaml``` as your expect.

The rules of this file are down below the last section.

### 3. generate running scripts
Run following command:

```
go build -o ns
./ns
```
or

```
go run main.go
```
After that, there should be some more shell scripts file in root folder:

```
launch.sh
clean.sh
```
### 4. start simulate network
Run ```launch.sh```

Then all thunderdb testnet docker should be running, and a graph drawed base on this network is lying in the root folder, which name ```graph.png```.

### 5. run your own program in testnet
The containers are name after their group_name+ip. For example, there are containers named 10.1.1.2 and 10.8.1.2, you can running ```docker exec -it china10.1.1.2 ping 10.8.1.2``` to test network connection between these two networks.

Replace "ping 10.8.1.2" to any program or script you like.

### 6. [optional]clean network
Run ```clean.sh```

## Modify Network Definition
The network description sample is in ```example.yaml```. You can edit it directly.

### sample
```
group:
  -
    name: china
    nodes:
        - 10.1.1.2/24
        - 10.2.1.1/16
        - 10.3.1.1/20
        - 10.4.1.1/20
    delay: "100ms 10ms 30%"
    loss: "1% 10%"
  -
    name: eu
    nodes:
        - 10.5.1.1/20
        - 10.6.1.1/20
        - 10.7.1.1/20
    delay: "10ms 5ms 30%"
    loss: "1% 10%"
  -
    name: jpn
    nodes:
        - 10.8.1.2/24
        - 10.9.1.2/24
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

## Description
The definition contains two section: group and network. Group defines ips and describes network info between them. Network describes network info between groups.

### group
- **name**: unique name of this group

- **node**: list ips of this network. Can only between "8.x.x.2 ~ 15.x.x.254" and must written in CIDR format.

- **network params**:
Support 6 tc network simulate params:

    ```
    delay
    loss
    duplicate
    corrupt
    reorder
    rate
    ```
The value of these params are exactly like ```tc``` command.

### network
- **groups**: list of group names.

- **network params**: same as group section
