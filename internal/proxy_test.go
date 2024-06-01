package internal

import (
    "testing"
)

func TestProxy(t *testing.T) {
    // UDP Server
    s_conn, err := StartServer("0.0.0.0", "8080")
    if err != nil {
	t.Error(err)
    }
    // UDP Client
    c_conn, err := StartClient("localhost", "8080")
    if err != nil {
	t.Error(err)
    }

    defer c_conn.Close()
    defer s_conn.Close()

    go func() {
	var buf [1024]byte
	for {
	    // Read from server
	    n, err := s_conn.Read(buf[:])
	    if err != nil {
		t.Error(err)
	    }
	    if string(buf[:n]) == "Hello, World!" {
		t.Log("Received: ", string(buf[:n]))
	    } else {
		t.Error("Received: ", string(buf[:n]))
	    }
	    return
	} 
    }()

    c_conn.Write([]byte("Hello, World!"))
    t.Log("Sent: Hello, World!")
}
