package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
)

type config struct {
	addr      uint64
	staticDir string
	logToFile bool
}

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
)

func main() {
	var cfg config

	flag.Uint64Var(&cfg.addr, "addr", 4000, "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.BoolVar(&cfg.logToFile, "log", false, "Enable logging")
	flag.Parse()

	InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if cfg.logToFile {
		infoFile, err := os.OpenFile("tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		
		err = infoFile.Truncate(0)
		if err != nil {
			log.Fatal(err)
		}

		errFile, err := os.OpenFile("tmp/error.log", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}

		err = errFile.Truncate(0)
		if err != nil {
			log.Fatal(err)
		}

		defer infoFile.Close()
		defer errFile.Close()
		InfoLog = log.New(infoFile, "INFO\t", log.Ldate|log.Ltime)
		ErrorLog = log.New(errFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	srv := &http.Server{
		Addr:     ":" + strconv.FormatUint(cfg.addr, 10),
		ErrorLog: ErrorLog,
		Handler:  mux,
	}

	InfoLog.Printf("Starting server on http://localhost:%d", cfg.addr)
	err := srv.ListenAndServe()
	ErrorLog.Fatal(err)
}
