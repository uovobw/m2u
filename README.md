m2u
===

Multicast-to-unicast converter in go

Small utility program i wrote to listen to a given multicast address and port and to
unicast all received data to many destinations

Not better than socat, a bit easier to use

Usage:

m2u multicastListenAddress:port unicastDestination:port unicastDestination:port ...
