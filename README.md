## fvpn

fvpn is a modern smart vpn. using fvpn we can create our own private network, as our devices such as cell phone, home pc and so on can join in this network, 
so they can connect to each other.

## Get Started
you can start fvpn from an example below:
-  register on our website: www.efvpn.com, after you registerd, you can free create three networks, if you want create more, you can see here.
-  install fvpn on your devices

## Galary
- Node: a node is a device which install in your pc/phone/...
- Registry: a registry is a service center, which node registered to
- NetworkId: a network is exactly a cidr, can be set in pc router.

## security
- DH
- curve-25519
- protocol-noise
beside of, every tunnel have different encrypt key.

## use insight
- working at home
- visit you nas when you are not at home
- visit your on private network when you don not want to buy a machine has public ip.

## Example
when you use fvpn, create an app.yaml is optional, because fvpn have a default config. when you create you own registry, you can change default config from create an app.yaml.
you should put app.yaml to /etc/fvpn or ~/.fvpn/
```azure
client:
  listen: :3000
  server: 127.0.0.1

#-------------------分害线
server:
  listen: :4000
  httpListen: :4009

```

start a node use commands below:
```shell
fvpn node
```

join a network: 
provide a networkId which you create on console.
```shell
fvpn join xxx
```

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

