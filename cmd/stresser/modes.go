package main

import (
    "context"
    "log"
    "sync"
)

type Stats struct {
    unmatched uint32
    matched uint32
}

type DNSTask struct {
    domain string 
}

type DNSTaskResult struct {
    domain string
    udguard_ip string
    cloudflare_ip string
} 

var wg sync.WaitGroup

func Worker(tasks chan DNSTask, results chan DNSTaskResult) {
    defer wg.Done()
    for task := range tasks {
	r_udguard := createDNSQuery(task.domain, UDGUARD_RESOLVER)
	udguard_ip, err := r_udguard.LookupIPAddr(context.Background(), task.domain)
	if err != nil {
	    logger.Printf("Failed to resolve domain %s: %v", task.domain, err)
	}

	r_cloudflare := createDNSQuery(task.domain, CLOUDFLARE_RESOLVER)
	cloudflare_ip, err := r_cloudflare.LookupIPAddr(context.Background(), task.domain)
	if err != nil {
	    logger.Printf("Failed to resolve domain %s: %v", task.domain, err)
	}

	var udg_ip, clf_ip string
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
	results <- DNSTaskResult{task.domain, udg_ip, clf_ip}
    }
}

func MultiThreaded() {
    domains := loadDomains()
    var NUM_TASKS = ITERATIONS

    tasks := make(chan DNSTask, NUM_TASKS)
    results := make(chan DNSTaskResult, NUM_TASKS)

    // Start workers
    for i := 0; i < NUM_WORKERS; i++ {
	wg.Add(1)
	go Worker(tasks, results)
    }

    // Send tasks to workers
    for i := 0; i < NUM_TASKS; i++ { 
	tasks <- DNSTask{domains[i]}
    }

    close(tasks)

    // Wait for workers to finish
    wg.Wait()

    // Collect results
    var stats Stats = Stats{0, 0}
    for i := 0; i < NUM_TASKS; i++ {
	result := <-results
	if result.udguard_ip != result.cloudflare_ip {
	    log.Printf("❌Domain: %s -> Cloudflare: %s, UDGuard: %s\n",result.domain, result.udguard_ip, result.cloudflare_ip)
	    stats.unmatched++
	} else {
	    log.Printf("✅Domain %s -> Cloudflare: %s, UDGuard: %s\n", result.domain, result.udguard_ip, result.cloudflare_ip)
	    stats.matched++
	}
    }

    log.Printf("Matched: %d, Unmatched: %d\n", stats.matched, stats.unmatched)
}

func SingleThreaded() {
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
