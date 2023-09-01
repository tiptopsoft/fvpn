## Relay

Relay used as a registry center, we support a Relay server default. Bug you can build your own relay if you want.

## What Relay do?

- node register to Relay server, Relay will cache Node info format: [ip: natIP:natPort], ip get from fvpn0 device in
  linux and windows, but on macos name like 'utunxx'
- node get other nodes from Relay every 30 secs, can build p2p to other nodes.

## When you should build your own Relay server?

- if node connect to our Relay slow, you need more fast register

## How wo build own Relay ?

- first you should have a public IP, and a cloud virtual machine near to you
- use simple command below

```shell
fvpn registry
```

will create a Relay server your self

## How to use own Relay server ?

- if you have a domain, you can set it on dns, if not, you can use ip directly

```shell
node:
  registry: your.domain:4000
```

port: 4000 is optional, if you leave it empty, default 4000 in code will be used.