package main

import (
	"flag"
	"os"
	"runtime"

	"github.com/itang/gohttp"
)

var (
	port    = 8080
	webroot = "."
)

func init() {
	wd, _ := os.Getwd()
	flag.IntVar(&port, "port", port, "The port (default is 8080)")
	flag.StringVar(&webroot, "webroot", wd, "Web root directory (default is current work directory)")

	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	server := &gohttp.FileServer{Port: port, Webroot: webroot}
	server.Start()
}
