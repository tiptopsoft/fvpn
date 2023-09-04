## Protocol Design

//TODO
We divide message five kind of type:

- HandShakeMessageType
- HandShakeMessageAckType
- MessagePacketType
- QueryNodeMessageType
- KeepaliveMessageType

## Packet Header

Every valid data message will add a header, header's length is 44, contains

```shell
Version uint8  //1
TTL     uint8  //1
Flags   uint16 //2
UserId  [8]byte
SrcIP   net.IP //16
DstIP   net.IP //16
```

- Version:  protocol version
- TTL: packet time to live
- Flag: message type, such as HandShakeMessageType and so on.
- UserId: user info, is a 10bit data
- SrcIP: ip of source node
- DstIP: dest of node
