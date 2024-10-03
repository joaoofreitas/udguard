# UDGuard

UDGuard is a base implementation of a zero-dependency DNS Proxy that enables authentication and packet processing. The data can be altered by adding custom code in the `DNSLookupHandler` function.

## Function to Modify
To customize the DNS Proxy behavior, modify the following function in `cmd/udguard/main.go`:

```go
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

    // ADD CODE HERE <------------

    log.Println("Sending response to client")
    _, err = s_conn.WriteToUDP(buf[:], addr)
    if err != nil {
        logger.Fatal(err)
    }
}
```

## Main Programs
There are two main programs in this repository located under the `cmd/` folder:
- `udguard`: The main DNS Proxy server.
- `stresser`: A tool for stress testing the proxy.

## Building and Compiling
To build and compile the programs, follow these steps:

1. Clone the repository:
   ```sh
   git clone https://github.com/joaoofreitas/udguard.git
   cd udguard
   ```

2. Build the main `udguard` program:
   ```sh
   cd cmd/udguard
   go build -o udguard
   ```

3. Build the `stresser` program:
   ```sh
   cd cmd/stresser
   go build -o stresser
   ```

## Testing
To test the `udguard` DNS Proxy, you can use the `stresser` tool:

1. Start the `udguard` server:
   ```sh
   ./udguard
   ```

2. In another terminal, run the `stresser` tool:
   ```sh
   ./stresser
   ```

The `stresser` tool will send multiple DNS requests to the `udguard` server to test its performance and stability.

## Responsibility
This code is provided as a base for implementing a DNS Proxy and is intended for educational and testing purposes only. It is not written with bad intentions, and I am not responsible for any misuse of this code, including its use as a basis for malware.
