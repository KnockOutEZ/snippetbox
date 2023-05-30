package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	conf "github.com/nexentra/snippetbox/configs"
	"github.com/nexentra/snippetbox/internal/models"
)

type config struct {
	addr      uint64
	staticDir string
	logToFile bool
}

type Application conf.Application

func main() {
	var cfg config
	var errorLog *log.Logger
	var infoLog *log.Logger

	dsn := flag.String("dsn", "postgres://postgres:mysecretpassword@localhost:5433/snippetbox?sslmode=disable", "Postgres data source name")

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

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &Application{
		ErrorLog:       errorLog,
		InfoLog:        infoLog,
		Snippets:       &models.SnippetModel{DB: db},
		TemplateCache:  templateCache,
		FormDecoder:    formDecoder,
		SessionManager: sessionManager,
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:     ":" + strconv.FormatUint(cfg.addr, 10),
		ErrorLog: app.ErrorLog,
		Handler:  app.routes(mux),
	}

	app.InfoLog.Printf("Starting server on http://localhost:%d", cfg.addr)
	err = srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
