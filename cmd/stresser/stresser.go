package main

import (
    "log"
    "os"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

const UDGUARD_RESOLVER = "127.0.0.1:8080"
const CLOUDFLARE_RESOLVER = "1.1.1.1:53"
const ITERATIONS = 1000
const NUM_WORKERS = ITERATIONS / 2

func main() {
    // Get OS Args
    args := os.Args[1:]
    if len(args) != 1 {
	logger.Printf("Usage: %s <mode>", os.Args[0])
	logger.Fatalf("Mode can be either 'single' or 'multi'")
    }

    mode := args[0]
    if mode == "single" {
	SingleThreaded()
    } else {
	MultiThreaded()
    }
}
