<img src="logo/logo.jpeg" width=200>

# GNTE
GNTE(Global Network Topology Emulator) is a docker-based all-in-one unstable global network emulator. It emulates functionality such as random delay and packet loss.

## Before Use
Install docker

## Build and Run
### 1. build docker image
Clone this repo and run ```build.sh```. There should be an image named ```ns``` in your docker environment.

### 2. modify network definition file
Edit ```example.yaml``` to fit your requirements. The rules of this file are described in the bottom section.

### 3. generate running scripts
Run the following command:

```
go build -o ns
./ns
```
or

```
go run main.go
```
Afterwards, your root folder should contain two shell scripts:

```
launch.sh
clean.sh
```
### 4. launch network emulator
Run ```launch.sh```

Once all thunderdb testnet dockers are running, you can use ```docker ps -a``` to see all container nodes: 
<img src="logo/container_node.png">

You can also find a graph of the network in ```graph.png``` under your root folder:
<img src="logo/graph.png">

### 5. run your own program in testnet
Containers are referenced by group_name+ip. For example, given containers 10.1.1.2 and 10.8.1.2, you can run ```docker exec -it china10.1.1.2 ping 10.8.1.2``` to test the connection between these two networks.

You can replace "ping 10.8.1.2" in the example above with any program or script.

### 6. [optional] clean network
Run ```clean.sh```

## Modify Network Definition
A sample network description is provided in ```example.yaml```, which you can edit directly.

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
The network definition contains two sections: group and network. Group defines ips and describes network info between them. Network describes network info between groups.

### group
- **name**: unique name of the group

- **node**: list of ips in the network. Must be between "8.x.x.2 ~ 15.x.x.254" and written in CIDR format

- **network params**:
The following 6 tc network limit parameters are supported:
    ```
    delay
    loss
    duplicate
    corrupt
    reorder
    rate
    ```
The values of these parameters are exactly like those of the ```tc``` command.

### network
- **groups**: list of group names

- **network params**: same as group
