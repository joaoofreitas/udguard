package internal

import (
    "net"
)

func StartServer(ip string, port string) (*net.UDPConn, error) {
    udpAddr, err := net.ResolveUDPAddr("udp", ip + ":" + port)
    if err != nil {
	return nil, err
    }
    conn, err := net.ListenUDP("udp", udpAddr)
    return conn, nil
}

func StartClient(ip string, port string) (*net.UDPConn, error) {
    udpAddr, err := net.ResolveUDPAddr("udp", ip + ":" + port)
    if err != nil {
	return nil, err
    }
    conn, err := net.DialUDP("udp", nil, udpAddr)
    return conn, nil
}
