package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

type config struct {
	addr      uint64
	staticDir string
}

func main() {
	var cfg config
	flag.Uint64Var(&cfg.addr, "addr", 4000, "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	log.Printf("Starting server on http://localhost:%d", cfg.addr)
	err := http.ListenAndServe(":"+strconv.FormatUint(cfg.addr, 10), mux)
	log.Fatal(err)
}
