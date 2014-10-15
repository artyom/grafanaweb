package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/artyom/autoflags"
	"github.com/artyom/basicauth"
)

func main() {
	config := struct {
		Addr   string `flag:"listen,address to listen"`
		Root   string `flag:"root,path to grafana files"`
		Prefix string `flag:"prefix,url prefix to db (should finish with /)"`
		Proxy  string `flag:"proxy,url to proxy prefix matched requests"`
		Auth   bool   `flag:"auth,use basic authentication (needs -authfile to be set)"`
		Creds  string `flag:"authfile,path to file with \"username:brypt_hash\" credential records"`

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
		log.Fatal("when using -ssl option, both -key and -cert options should also be set")
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

	switch {
	case config.Auth && len(config.Creds) == 0:
		log.Fatal("when using -auth option, -authfile options should also be set")
	case config.Auth && len(config.Creds) > 0:
		realm := basicauth.NewRealm("Restricted area")
		if err := loadCredentials(realm, config.Creds); err != nil {
			log.Fatal(err)
		}
		//realm.AddUser("Aladdin", "open sesame")
		http.Handle(config.Prefix, realm.WrapHandler(httputil.NewSingleHostReverseProxy(url)))
		http.Handle("/", realm.WrapHandler(http.FileServer(http.Dir(config.Root))))
	default:
		http.Handle(config.Prefix, httputil.NewSingleHostReverseProxy(url))
		http.Handle("/", http.FileServer(http.Dir(config.Root)))
	}

	server := &http.Server{
		Addr:           config.Addr,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
		TLSConfig:      &tls.Config{MinVersion: tls.VersionTLS10},
	}

	if config.SSL {
		log.Fatal(server.ListenAndServeTLS(config.Cert, config.Key))
	}
	log.Fatal(server.ListenAndServe())
}

// loadCredentials reads credentials from file and loads them to realm. Empty
// lines or lines beginning with # (exactly first characted) are ignored, other
// lines are treated as colon-separated username and password bcrypt hash. Use
// `bcryptpasswd` utility to produce such records.
func loadCredentials(realm *basicauth.Realm, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var hashedSecret []byte
	var n int
	for scanner.Scan() {
		n++
		// skip empty lines and comments
		if len(scanner.Bytes()) == 0 || scanner.Bytes()[0] == '#' {
			continue
		}
		s := strings.SplitN(scanner.Text(), ":", 2)
		if len(s) != 2 {
			return fmt.Errorf("invalid record in %q line %d", path, n)
		}
		username := s[0]
		hashedSecret, err = base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			return fmt.Errorf("invalid record in %q line %d: %s", path, n, err)
		}
		if err := realm.AddUserHashed(username, hashedSecret); err != nil {
			return fmt.Errorf("invalid record in %q line %d: %s", path, n, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
