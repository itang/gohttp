package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var (
	PORT    = 8080
	WEBROOT = "."
)

type Server struct {
	Port    int
	Webroot string
}

func init() {
	wd, _ := os.Getwd()
	log.Printf("current dir:%v", wd)
	log.Printf("PathSeparator:%c", os.PathSeparator)

	flag.IntVar(&PORT, "port", PORT, "the port (default is 8080)")
	flag.StringVar(&WEBROOT, "webroot", wd, "Web root directory (default is current work directory)")
}

func (server *Server) router() {
	http.Handle("/", server)
}

func (server *Server) Start() {
	log.Printf("port:%v", server.Port)
	log.Printf("webroot:%v", server.Webroot)

	server.router()

	addr := fmt.Sprintf(":%v", server.Port)
	fmt.Printf("start server at :%v\n", server.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Errorf("%v", err)
	}
}

func (server *Server) requestURIToFilepath(uri string) (fullpath string, relpath string) {
	unescapeIt, _ := url.QueryUnescape(uri)

	relpath = unescapeIt
	log.Printf("unescape uri:%v", relpath)

	fullpath = filepath.Join(server.Webroot, relpath[1:])
	//** trace
	log.Printf("base path:%v\n", path.Base(fullpath))
	log.Printf("dir path:%v\n", path.Dir(fullpath))
	log.Printf("ext path:%v\n", path.Ext(fullpath))
	//**
	return
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type Item struct {
	Name  string
	Title string
	URI   string
	Size  int64
}

func (server *Server) processDir(w http.ResponseWriter, dir *os.File, fullpath string, relpath string) {
	w.Header().Set("Content-type", "text/html; charset=UTF-8")
	fis, err := dir.Readdir(-1)
	checkError(err)

	items := make([]Item, 0, len(fis))
	for _, fi := range fis {
		item := Item{
			Name:  fi.Name(),
			Title: fi.Name(),
			URI:   path.Join(relpath, fi.Name()),
			Size:  fi.Size(),
		}
		items = append(items, item)
	}

	tmp.Execute(w, map[string]interface{}{
		"ParentURI":  path.Dir(relpath),
		"CurrentURI": relpath,
		"files":      items,
	})
}

func (server *Server) sendFile(w http.ResponseWriter, file *os.File, fullpath string, relpath string) {
	if mimetype := mime.TypeByExtension(path.Ext(file.Name())); mimetype != "" {
		w.Header().Set("Content-Type", mimetype)
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	statinfo, _ := file.Stat()
	w.Header().Set("Content-Length", fmt.Sprintf("%v", statinfo.Size()))
	io.Copy(w, file)
}

func (server *Server) handler(w http.ResponseWriter, req *http.Request) {
	uri := req.RequestURI      // 请求的URI, 如http://localhost:8080/hello -> /hello
	if uri == "/favicon.ico" { // 不处理
		return
	}

	fullpath, relpath := server.requestURIToFilepath(uri)
	log.Printf("To Filepath:%v", fullpath)

	file, err := os.Open(fullpath)
	if err != nil || os.IsNotExist(err) { // 文件不存在
		http.NotFound(w, req)
	} else {
		stat, _ := file.Stat()

		log.Printf("%v is dir? %v\n", file.Name(), stat.IsDir())
		if stat.IsDir() {
			server.processDir(w, file, fullpath, relpath)
		} else {
			server.sendFile(w, file, fullpath, relpath)
		}
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("<< Request from %v", req.RemoteAddr)

	defer func() {
		if err := recover(); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		}
	}()

	server.handler(w, req)

	log.Printf(" End Request>>")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	server := &Server{Port: PORT, Webroot: WEBROOT}
	server.Start()
}

var tmp = template.Must(template.New("index").Parse(html))

const html = `
<a href="{{.ParentURI}}"> {{.ParentURI}} </a> | <a href="{{.CurrentURI}}">{{.CurrentURI}}</a>
<ul>
   {{range .files}}
      <li><a href="{{.URI}}">{{.Name}}
      {{if .Size }}
      <small>({{.Size}})</small>
      {{end}}
      </a></li>
   {{end}}
</ul>`
