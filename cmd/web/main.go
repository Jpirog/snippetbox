package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)
type application struct{
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
		infoLog: infoLog,
	}

	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// register the base path to the home function
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	fullname := fmt.Sprintf("%s:%d", *host, *port)
	srv := &http.Server{
		Addr:     fullname,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on port %v ", fullname)
	herr := srv.ListenAndServe()
	errorLog.Fatal(herr)
}
