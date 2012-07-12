package main

import (
	"fmt"
	"flag"
	"net/http"
	"net/url"
	"os"
	"path"
)

var (
	PORT    = "8080"
	WEBROOT = ""
)

func init() {
	flag.StringVar(&PORT, "port", PORT, "the port (default is 8080)")
	wd, _ := os.Getwd()
	flag.StringVar(&WEBROOT, "webroot", wd, "Web root directory (default is current work directory)")
}

func handle(w http.ResponseWriter, req *http.Request) {
  //w.Header().Set("Content-type", "text/html; charset=UTF-8")

	uri := req.RequestURI
	fmt.Fprintln(w, "hello", uri)

	un, _ := url.QueryUnescape(uri)
	fmt.Fprintln(w, "unescape uri:", un)

	p := path.Join(WEBROOT, un[1:])
	fmt.Fprintln(w, "path:", p)
	fmt.Fprintf(w, "base path:%v\n", path.Base(p))
	fmt.Fprintf(w, "dir path:%v\n", path.Dir(p))
	fmt.Fprintf(w, "ext path:%v\n", path.Ext(p))

	file, err := os.Open(p)
	if err != nil || os.IsNotExist(err) {
		http.NotFound(w, req)
	} else {
		fmt.Fprintf(w, "%v\n", file.Name())
		stat, _ := file.Stat()

		fmt.Fprintf(w, "%v is dir? %v\n", file.Name(), stat.IsDir())
		if stat.IsDir() {
			processDir(w, file)
		}

	}
}

func processDir(w http.ResponseWriter, dir *os.File) {
	fi, err := dir.Readdir(-1)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	for _, fi := range fi {
		//if fi.IsRegular() {
		fmt.Fprintln(w, fi.Name(), fi.Size(), "bytes")
		//}
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/", handle)

	fmt.Printf("start server at :%v\n", PORT)
	addr := fmt.Sprintf(":%v", PORT)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Errorf("%v", err)
	}
}
