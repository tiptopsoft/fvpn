## Fvpn

Fvpn is a smart vpn software inspired by wireguard. Its aim is allowing all kinds devices, like cell phones,vms,pcs and
so
on to communicate with each other, as if they all reside in the same physical data center or cloud region,
default each pair node will build p2p connection, all your node will compose a private network mesh.

## Quick start

- register on our website: https://www.tiptopsoft.cn
- install fvpn to your devices, follow this document: [install fvpn](./docs/install.md)
- using fvpn by a default network created for you, but which services is limited

## Glossary

- Node: node is a device which running in your vm/pc/phone/...
- Registry: registry is a data center, which node registered to, and also relay traffic when node behind Symmetric
  NAT
- NetworkId: network is exactly a cidr, can be set in pc/vm router.

## Why fvpn

as so much software you can choose, but fvpn has some differences:

- fvpn is simple & light, you can quickly compose your own private network, just follow a few steps.
- fvpn is quick, when another new device join in the network, other devices will immediately recognize and try to build
  a p2p to it.
- fvpn is safe, follow the protocol_noise & different DH key used in each pair node. and key can be changed manually.

## Usage scenarios

- working at home, but you want to visit your company pc or other devices
- visit your home nas anytime when wherever you are.
- visit your own private network anytime when you do not want to buy a machine has public ip.

## Configuration

when you use fvpn, create an app.yaml is optional, because fvpn have a default config. you just should modify
configuration when you create you own registry.
you can change default config from create an app.yaml.
you should put app.yaml to /etc/fvpn or ~/.fvpn/, then some configurations will be covered the default config.

```shell
node:
  registry: tiptopsoft.cn
```

## Build

First, clone code from our repository:

```shell
git clone git@github.com:tiptopsoft/fvpn.git
```

use make command:

```shell
cd fvpn && make build
```
will create exec file under bin folder.

## Contributions

All kind of contributions are welcome, issues,PR and so on.
