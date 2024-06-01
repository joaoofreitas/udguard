package main

import (
    "context"
    "time"
    "bufio"
    "log"
    "net"
    "os"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
    // Load domains from file
    absPath, _ := os.Getwd()
    logger.Printf("Current working directory: %s\n", absPath)

    file, err := os.Open("public-domain-lists/opendns-top-domains.txt")
    if err != nil {
	logger.Fatalf("Failed to open file: %v", err)
    }
    defer file.Close()

    var domains []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        domains = append(domains, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
	logger.Fatalf("Failed to read domains: %v", err)
    }
    logger.Printf("Loaded %d domains\n", len(domains))
}

func loadDomains() []string {
    var domains []string

    // Load domains from file
    file, err := os.Open("public-domain-lists/opendns-top-domains.txt")
    if err != nil {
	logger.Fatalf("Failed to open file: %v", err)
    }
    
    scanner := bufio.NewScanner(file)
    if err := scanner.Err(); err != nil {
	logger.Fatalf("Failed to read domains: %v", err)
    }
    logger.Printf("Loaded %d domains\n", len(domains))

    return domains
}

// Helper function to create a DNS query message
func createDNSQuery(domain string) *net.Resolver {
    r := &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
	    d := net.Dialer{
		Timeout: 5 * time.Second,
	    }
	    return d.DialContext(ctx, "udp", "1.1.1.1:53")
	},
    }
    return r
}

