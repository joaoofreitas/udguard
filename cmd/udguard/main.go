package main

import (
    "os"
    "log"
    "net"
    "github.com/joaoofreitas/udguard/internal"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
    s_conn, err := internal.StartServer("0.0.0.0", "8080")
    if err != nil {
	logger.Fatal(err)
	panic(err)
    }
    defer s_conn.Close()

    for {
	var buf [512]byte
	n, addr, err := s_conn.ReadFromUDP(buf[:])
	if err != nil {
	    logger.Fatal(err)
	    panic(err)
	}
	go DNSLookupHandler(buf[0:n], addr, s_conn)
    }
}

func DNSLookupHandler(msg []byte, addr *net.UDPAddr, s_conn *net.UDPConn) {
    c_conn, err := internal.StartClient("1.1.1.1", "53")
    if err != nil {
	logger.Fatal(err)
	panic(err)
    }
    defer c_conn.Close()

    log.Println("Sending request to DNS server")
    _, err = c_conn.Write(msg)

    var buf [512]byte
    _, _, err = c_conn.ReadFromUDP(buf[0:])
    if err != nil {
	logger.Fatal(err)
    }
    log.Println("Received response from DNS server")
    log.Println(buf)
    
    log.Println("Waiting for response from DNS server")

    log.Println("Sending response to client")
    _, err = s_conn.WriteToUDP(buf[:], addr)
    if err != nil {
	logger.Fatal(err)
    }
}
