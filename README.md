## fvpn

A fvpn is a modern smart vpn. use fvpn we can visit our home network when we are outside.

## Get Started
fvpn have client and server, you can use fvpn compose you own network. fvpn need at least one public Node which has an public IP.

in public node:
```shell
fvpn -s -pxxx
```

in our private node:
```shell
fvpn -c xx.xx.xx.xx -pxxx
```

we must  use the same port on our server and client. if you did not specify one, the default port 6061 will be used.

## Compile

First, clone code from our repository:

```shell
git clone
```

use make command:

```shell
cd fvpn && make build
```

## Contribute

- closed
