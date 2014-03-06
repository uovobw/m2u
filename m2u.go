package main

import (
    "fmt"
    "os"
    "net"
    )

var listenAddr string
var sendAddrOne string
var sendAddrTwo string
var addrToChan map[string]chan []byte

func unicast_send(datain chan([]byte) , addr string) {
    fmt.Print("Starting for address: ")
    fmt.Println(addr)
    address, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
        fmt.Println("Cannot resolve passed address")
        os.Exit(1)
    }
    conn, err := net.DialUDP("udp", nil, address)
    if err != nil {
        fmt.Print("Cannot ListenUDP to addr: ")
        fmt.Println(err)
        os.Exit(1)
    }
    // start reading from the chan and sending to the remote host
    for {
        d := <- datain
        _, err := conn.Write(d)
        if err != nil {
            // for the time being, do nothing, we might still not have someone
            // listening on the other side
        }
    }
}

func main() {

    addrToChan = make(map[string]chan []byte)

    if len(os.Args) < 4 {
        fmt.Println("Usage: m2u multicastGroup:port unicastDestination:port unicastDestination:port [unicastDestination:port ...]")
        os.Exit(1)
    }

    listenAddr = os.Args[1]

    for _, h := range os.Args[2:] {
        addrToChan[h] = make(chan []byte, 2048)
    }

    fmt.Println("Starting")
    mcaddr, err := net.ResolveUDPAddr("udp", listenAddr)
    if err != nil {
        fmt.Println("Could not resolve multicast address")
        os.Exit(1)
    }
    socket, err := net.ListenMulticastUDP("udp4", nil, mcaddr)
    if err != nil {
        fmt.Println("Could not ListenMulticastUDP")
        os.Exit(1)
    }
    defer socket.Close()

    // spawn as many goroutines ad needed
    for addr, channel := range addrToChan {
        go unicast_send(channel, addr)
    }

    // start reading
    for {
        b := make([]byte, 1500)
        n, _, err := socket.ReadFromUDP(b)
        if err != nil {
            fmt.Println("Could not ReadFromUDP")
        }
        toWrite := b[:n]
        for _, c := range addrToChan {
            c <- toWrite
        }
    }

}
