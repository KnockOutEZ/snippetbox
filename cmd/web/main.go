package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	conf "github.com/nexentra/snippetbox/configs"
)

type config struct {
	addr      uint64
	staticDir string
	logToFile bool
}

type Logger conf.Logger

func main() {
	var cfg config
	var errorLog *log.Logger
	var infoLog *log.Logger

	flag.Uint64Var(&cfg.addr, "addr", 4000, "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.BoolVar(&cfg.logToFile, "log", false, "Enable logging")
	flag.Parse()

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
		infoLog = log.New(infoFile, "INFO\t", log.Ldate|log.Ltime)
		errorLog = log.New(errFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	}

	app := &Logger{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	srv := &http.Server{
		Addr:     ":" + strconv.FormatUint(cfg.addr, 10),
		ErrorLog: app.ErrorLog,
		Handler:  mux,
	}

	app.InfoLog.Printf("Starting server on http://localhost:%d", cfg.addr)
	err := srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}
