package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/qbitty/snippetbox/pkg/config"
	"github.com/qbitty/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog *log.Logger
	errLog  *log.Logger
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	dns := flag.String("dns", "web:pass@/snippetbox?parseTime=true", "Mysql database dns")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dns)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	app := &config.Application{
		ErrLog:   errLog,
		InfoLog:  infoLog,
		Snippets: &mysql.SnippetModel{DB: db},
	}

	// Use the http.ListenAndServe() function to start a new web server. We pas
	// two parameters: the TCP network address to listen on (in this case ":4000
	// and the servemux we just created. If http.ListenAndServe() returns an er
	// we use the log.Fatal() function to log the error message and exit.
	// infoLog.Printf("Starting server on %s", *addr)
	// err := http.ListenAndServe(*addr, mux)
	// errLog.Fatal(err)

	svr := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  routes(app),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = svr.ListenAndServe()
	errLog.Fatal(err)
}

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
