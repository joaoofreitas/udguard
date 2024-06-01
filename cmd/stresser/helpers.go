package main

import (
    "context"
    "time"
    "bufio"
    "os"
    "net"
)

func loadDomains() []string {
    var domains []string
    // Load domains from file
    file, err := os.Open("public-domain-lists/opendns-top-domains.txt")
    if err != nil {
	logger.Fatalf("Failed to open file: %v", err)
    }
    
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
	domains = append(domains, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
	logger.Fatalf("Failed to read domains: %v", err)
    }
    logger.Printf("Loaded %d domains\n", len(domains))

    return domains
}

// Helper function to create a DNS query message
func createDNSQuery(domain string, resolver string) *net.Resolver {
    r := &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
	    d := net.Dialer{
		Timeout: 5 * time.Second,
	    }
	    return d.DialContext(ctx, "udp", resolver)
	},
    }
    return r
}
