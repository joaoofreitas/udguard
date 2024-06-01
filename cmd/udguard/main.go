package main

import (
    "os"
    "log"
    "github.com/joaoofreitas/udguard/internal"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
    s_conn, err := internal.StartServer("0.0.0.0", "8080")
    if err != nil {
	logger.Fatal(err)
	panic(err)
    }
    
    c_conn, err := internal.StartClient("0.0.0.0", "8081")
    if err != nil {
	logger.Fatal(err)
	panic(err)
    }
    
    c_to_s:= make(chan []byte)
    s_to_c := make(chan []byte)
    go func(c_to_s chan []byte) {
	var buf [1024]byte
	for {
	    n, err := s_conn.Read(buf[:])
	    if err != nil {
		logger.Fatal(err)
		panic(err)
	    }
	    c_to_s <- buf[:n]
	}
    }(c_to_s)

    go func(s_to_c chan []byte) {
	var buf [1024]byte
	for {
	    n, err := c_conn.Read(buf[:])
	    if err != nil {
		logger.Fatal(err)
		panic(err)
	    }
	    s_to_c <- buf[:n]
	}
    }(s_to_c)

    for {
	select {
	case msg := <- c_to_s:
	    _, err := s_conn.WriteTo(msg, c_conn.RemoteAddr()) 

	    if err != nil {
		logger.Fatal(err)
	    }
	    logger.Println("Sent: ", string(msg))
	case msg := <- s_to_c:
	    _, err := c_conn.WriteTo(msg, c_conn.RemoteAddr()) 

	    if err != nil {
		logger.Fatal(err)
	    }
	    logger.Println("Received: ", string(msg))
	}
    }
}
