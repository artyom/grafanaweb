package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/artyom/autoflags"
)

func main() {
	config := struct {
		Addr   string `flag:"listen,address to listen"`
		Root   string `flag:"root,path to grafana files"`
		Prefix string `flag:"prefix,url prefix to db (should finish with /)"`
		Proxy  string `flag:"proxy,url to proxy prefix matched requests"`

		SSL  bool   `flag:"ssl,use https instead of http (needs both -cert and -key options set)"`
		Key  string `flag:"key,path to key.pem file (only if -ssl used)"`
		Cert string `flag:"cert,path to cert.pem file (only if -ssl used)"`
	}{
		Addr:   "127.0.0.1:8080",
		Root:   "/var/lib/grafana",
		Prefix: "/db/",
		Proxy:  "http://127.0.0.1:8086",
	}
	if err := autoflags.Define(&config); err != nil {
		log.Fatal(err)
	}
	flag.Parse()

	if config.SSL && (len(config.Key) == 0 || len(config.Cert) == 0) {
		log.Fatal("when using -ssl option, both -key and -cert options should be set")
	}
	if !strings.HasSuffix(config.Prefix, "/") {
		config.Prefix = config.Prefix + "/"
	}
	url, err := url.Parse(config.Proxy)
	if err != nil {
		log.Fatal(err)
	}
	if !url.IsAbs() {
		log.Fatal("proxy url should be absolute")
	}
	http.Handle(config.Prefix, httputil.NewSingleHostReverseProxy(url))
	http.Handle("/", http.FileServer(http.Dir(config.Root)))
	if config.SSL {
		log.Fatal(http.ListenAndServeTLS(config.Addr, config.Cert, config.Key, nil))
	}
	log.Fatal(http.ListenAndServe(config.Addr, nil))
}
