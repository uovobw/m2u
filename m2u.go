package main

import (
    "fmt"
    "os"
    "net"
    "flag"
    )

var listenAddr string
var addrToChan map[string]chan []byte
var verbose *bool

func unicast_send(datain chan([]byte) , addr string) {
    if *verbose {
        fmt.Println(fmt.Sprintf("Starting unicast sender to %s", addr))
    }
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
    defer conn.Close()
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

    verbose = flag.Bool("v", false, "be verbose")
    flag.Parse()

    addrToChan = make(map[string]chan []byte)

    if len(os.Args) < 4 {
        fmt.Println(fmt.Sprintf("Usage: %s [-v] multicastGroup:port unicastDestination:port [unicastDestination:port ...]", os.Args[0]))
        os.Exit(1)
    }

    if *verbose {
        listenAddr = os.Args[2]
        for _, h := range os.Args[3:] {
            addrToChan[h] = make(chan []byte, 2048)
        }
    } else {
        listenAddr = os.Args[1]
        for _, h := range os.Args[2:] {
            addrToChan[h] = make(chan []byte, 2048)
        }
    }

    if *verbose {
        fmt.Println(fmt.Sprintf("Multicast address: %s", listenAddr))
        for address, _ := range addrToChan {
            fmt.Println(fmt.Sprintf("Sending to: %s", address))
        }
    }

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

    // spawn as many goroutines as needed
    for addr, channel := range addrToChan {
        go unicast_send(channel, addr)
    }

    // main read from connection -> write to channel loop
    if *verbose {
        fmt.Println("Entering main loop")
    }
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
