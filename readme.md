gohttp: go simple http server
---------------------

Golang implementation replace "python -m SimpleHTTPServer"

### usage

```
$ go get -u github.com/itang/gohttp
$ go install github.com/itang/gohttp/gohttp
    
$ gohttp --help

$ gohttp
Serving HTTP on 192.168.1.103 port 8080 from "/home/itang/workspace/work" ... 

$ gohttp -d=/home -p=9000
Serving HTTP on 192.168.1.128 port 9000 from "/home" ...
