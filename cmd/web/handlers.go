package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// this is the base home handler
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file
	// paths as a variadic parameter?
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print("1" + err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print("2" + err.Error())
		http.Error(w, err.Error(), 500)
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
