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

var UDGUARD_RESOLVER = "127.0.0.1:8080"
var CLOUDFLARE_RESOLVER = "1.1.1.1:53"
var ITERATIONS = 1000

type Stats struct {
    unmatched uint32
    matched uint32
}


func main() {
    domains := loadDomains()
    var stats Stats = Stats{0, 0}

    // For each domain, create a DNS query, send it to 1.1.1.1 and 127.0.0.1 and compare the responses
    var udg_ip, clf_ip string
    
    for i, domain := range domains {
	r_udguard := createDNSQuery(domain, UDGUARD_RESOLVER)
	udguard_ip, err := r_udguard.LookupIPAddr(context.Background(), domain)
	if err != nil {
	    logger.Printf("Failed to resolve domain %s: %v", domain, err)
	}

	r_cloudflare := createDNSQuery(domain, CLOUDFLARE_RESOLVER)
	cloudflare_ip, err := r_cloudflare.LookupIPAddr(context.Background(), domain)
	if err != nil {
	    logger.Printf("Failed to resolve domain %s: %v", domain, err)
	}
	
	if len(udguard_ip) == 0 {
	    udg_ip = "NOT FOUND"
	} else {
	    udg_ip = udguard_ip[0].IP.String()
	}
	if len(cloudflare_ip) == 0 {
	    clf_ip = "NOT FOUND"
	} else {
	    clf_ip = cloudflare_ip[0].IP.String()
	}

	if udg_ip != clf_ip {
	    log.Printf("❌Domain: %s -> Cloudflare: %s, UDGuard: %s\n", domain,udg_ip, clf_ip)
	    stats.unmatched++
	} else {
	    log.Printf("✅Domain %s -> Cloudflare: %s, UDGuard: %s\n", domain, udg_ip, clf_ip)
	    stats.matched++
	}

	// Stop @ iterations
	if i == ITERATIONS {
	    log.Printf("Matched: %d, Unmatched: %d\n", stats.matched, stats.unmatched)
	    break
	}
    }
}

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

