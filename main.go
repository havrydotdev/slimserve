package main

import (
	"errors"
	"flag"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

var (
	port      = flag.String("p", "8100", "port to serve on")
	directory = flag.String("d", ".", "the directory of static file to host")
	encoding  = flag.String("e", "", "encoding to use")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := path.Join(*directory, r.URL.Path)

		if r.URL.Path == "/" {
			p = path.Join(*directory, "index.html")
		}

		w.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(p)))

		switch *encoding {
		case "brotli":
			if e, _ := exists(p + ".br"); e {
				w.Header().Add("Content-Encoding", "br")
				p = p + ".br"
			}
		case "gzip":
			if e, _ := exists(p + ".gz"); e {
				w.Header().Add("Content-Encoding", "gzip")
				p = p + ".gz"
			}
		}

		http.ServeFile(w, r, p)
	})

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err
}
