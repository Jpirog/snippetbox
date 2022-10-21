package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.johnpirog.com/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	host := flag.String("host", "localhost", "HTTP network address")
	port := flag.Int("port", 7999, "HTTP port")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Llongfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	dsn := flag.String("dsn", "root:mysql@tcp(localhost:3306)/snippetbox?parseTime=true", "MySQL data source name")
	infoLog.Printf("Connecting to DB: %v", *dsn)
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	fullname := fmt.Sprintf("%s:%d", *host, *port)
	srv := &http.Server{
		Addr:     fullname,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on port %v ", fullname)
	herr := srv.ListenAndServe()
	errorLog.Fatal(herr)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
