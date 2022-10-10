package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	host := flag.String("host", "localhost", "HTTP network address")
	port := flag.Int("port", 7999, "HTTP port")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Llongfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
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
