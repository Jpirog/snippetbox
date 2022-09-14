package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// this is the base home handler
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("foo", "bar")
	// w.Header().Set("ngrok-skip-browser-warning", "true")
	w.Header()["ngrok-skip-browser-warning"] = []string{"true"}
	// w.Header().Add("user-agent", "ngrok-go")
	// w.Header().Add("foo", "bar")

	// 	headers := http.Header{
	// 		"ngrok-skip-browser-warning": []string{"true"},
	// 		"Accept": []string{"text/plain", "text/html"},
	// }
	// r.Header["ngrok-skip-browser-warning"] = []string{"true"}

	log.Print("Responding to a request at ", time.Now())
	log.Print(" -- Host", r.Host, r.URL.Path)
	log.Print(" -- RemoteAddr", r.RemoteAddr)
	log.Print(w.Header())
	// w.WriteHeader(200)
	w.Write([]byte("Hello from Snippetbox, current timestamp is " + time.Now().String()))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	log.Print("snippetView request", r.URL.Path)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Displaying a specific snippet (%d)...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	log.Print("snippetCreate request", r.URL.Path)
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "Post")
		w.WriteHeader(405)
		w.Write(([]byte("Method not allowed")))
		return
	}
	w.Write([]byte("Create a new snippet..."))
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
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on port ", port)
	err = http.ListenAndServe(":"+port, mux)
	log.Fatal(err)
}
