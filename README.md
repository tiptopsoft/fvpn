## FVPN

A fvpn is a modern smart vpn. use fvpn we can create our own private network, after we join a network, we can connect other device that join in the network we just join in.
Use fvpn you also can compose you own registry, if you didn't, you can use our services online. as it is default choice.
## Get Started
you can start fvpn from an example below:

## Galary
- Node: a node is a device which install in your pc/router/...
- Registry: a registry is a service center, node can register its self to a registry
- NetworkId: 
- User: If you want to use our public services, you should sign up an account in our website. you can create networkId, you also can manage it. 

## Example
when you use fvpn, you should create an app.yml
```azure
client:
  listen: :3000
  server: 127.0.0.1
  tap: tap0
  ip: 192.168.1.1
  mask: 255.255.255.0
  mac: 01:02:0f:0E:04:01
  type: udp

#-------------------分害线
server:
  listen: :4000
  httpListen: :4009
  type: udp
```
start a node use commands below:
```azure
fvpn node
```
you can use alias short name:
```azure
fvpn n
```

join a network:
```azure
fvpn join xxx
```

## IP 
after you join a network, you can visit you console, in which you can see your device status and device info. also can see your devices register in this network. as a device can 
join multiple network.

## Compile

First, clone code from our repository:

```shell
git clone
```
## Build
use make command:

```shell
cd fvpn && make build
```

## Contribute

