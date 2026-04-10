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

./stresser -m https3 -t https://target.com -p proxies.txt -th 1000 -d 60
./stresser -m post -t https://target.com -th 800
./stresser -m https2 -t https://target.com -th 700

```
## L4 examples:
```bash

./stresser -m udppower -t 8.8.8.8:53 -th 1200 -d 30
./stresser -m tcpfull -t 192.168.1.1:80 -th 900
```
