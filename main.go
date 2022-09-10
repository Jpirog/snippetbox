package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// this is the base home handler
func home(w http.ResponseWriter, r *http.Request) {
	log.Print("Responding to a request at ", time.Now())
	log.Print(" -- Host", r.Host)
	log.Print(" -- RemoteAddr", r.RemoteAddr)
	w.Write([]byte("Hello from Snippetbox, current timestamp is " + time.Now().String()))
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	// register the base path to the home function
	mux.HandleFunc("/", home)

	log.Print("Starting server on port ", port)
	err = http.ListenAndServe(":"+port, mux)
	log.Fatal(err)
}
