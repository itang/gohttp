package main

import (
	"flag"
	"github.com/itang/gohttp"
	"os"
	"runtime"
)

var (
	port    = 8080
	webroot = "."
)

func init() {
	wd, _ := os.Getwd()
	flag.IntVar(&port, "port", port, "The port (default is 8080)")
	flag.StringVar(&webroot, "webroot", wd, "Web root directory (default is current work directory)")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()*2 - 1)
	flag.Parse()
	server := &gohttp.FileServer{Port: port, Webroot: webroot}
	server.Start()
}
